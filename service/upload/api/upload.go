package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kuan525/netdisk/client/dbproxy"
	cmn "github.com/kuan525/netdisk/common"
	cfg "github.com/kuan525/netdisk/config"
	"github.com/kuan525/netdisk/mq"
	"github.com/kuan525/netdisk/service/dbproxy/orm"
	"github.com/kuan525/netdisk/store/ceph"
	"github.com/kuan525/netdisk/store/cos"
	"github.com/kuan525/netdisk/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// DoUploadHandler 处理文件上传
func DoUploadHandler(c *gin.Context) {
	errCode := 0

	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传成功",
			})
		}
	}()
	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()

	// 1. 从form表单中获得文件内容句柄
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to get form data, err:%s\n", err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	// 2. 把文件内容转为[]byte
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Printf("Failed to get file data, err:%s\n", err.Error())
		errCode = -2
		return
	}

	// TODO 3. 构建文件元信息
	fileMeta := dbproxy.FileMeta{
		FileName: head.Filename,
		FileSha1: util.Sha1(buf.Bytes()), // 计算文件sha1
		FileSize: int64(len(buf.Bytes())),
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 4. 将文件写入临时存储位置
	fileMeta.Location = cfg.TempLocalRootDir + fileMeta.FileSha1 // 临时存储地址
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		log.Printf("Failed to create file, err:%s\n", err.Error())
		errCode = -3
		return
	}
	defer newFile.Close()

	nByte, err := newFile.Write(buf.Bytes())
	if int64(nByte) != fileMeta.FileSize || err != nil {
		log.Printf("Failed to save data into file, writtenSize:%d, err:%s\n", nByte, err.Error())
		errCode = -4
		return
	}

	// 5. 同步或异步将文件转移到Ceph/COS
	newFile.Seek(0, 0) // 游标重新回到文件头部
	if cfg.CurrentStoreType == cmn.StoreCeph {
		// 文件写入Ceph存储
		data, _ := ioutil.ReadAll(newFile)
		cephPath := cfg.CephRootDir + fileMeta.FileSha1
		_ = ceph.PutObject("userfile", cephPath, data)
		fileMeta.Location = cephPath
	} else if cfg.CurrentStoreType == cmn.StoreCOS {
		// 文件写入cos存储
		cosPath := cfg.COSRootDir + fileMeta.FileSha1
		// 判断写入cos为同步还是异步
		if !cfg.AsyncTransferEnable {
			// TODO 设置cos中的文件名，方便指定文件名下载
			//cos.NewClient().B
			_, err := cos.Client().Object.Put(context.Background(), cosPath, newFile, nil)
			if err != nil {
				log.Println(err.Error())
				errCode = -5
				return
			}
			fileMeta.Location = cosPath
		} else {
			// 写入异步转移任务队列
			data := mq.TransferData{
				FileHash:      fileMeta.FileSha1,
				CurLocation:   fileMeta.Location,
				DestLocation:  cosPath,
				DestStoreType: cmn.StoreCOS,
			}
			pubData, _ := json.Marshal(data)
			puSuc := mq.Publish(
				cfg.TranExchangeName,
				cfg.TransCOSRoutingKey,
				pubData,
			)
			if !puSuc {
				// TODO 当前发送转移信息失败，稍后重试
			}
		}
	}

	// 6. 更新文件表记录
	_, err = dbClient.OnFileUploadFinished(fileMeta)
	if err != nil {
		errCode = -6
		return
	}

	// 7. 更新用户文件表
	username := c.Request.FormValue("username")
	upRes, err := dbClient.OnUserFileUploadFinished(username, fileMeta)
	if err != nil && upRes.Suc {
		errCode = 0
	} else {
		errCode = -6
	}
}

// TryFastUploadHandler 尝试秒传接口
func TryFastUploadHandler(c *gin.Context) {
	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()

	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")

	// 2. 从文件表中查询相同的hash的文件记录
	fileMetaResp, err := dbClient.GetFileMeta(filehash)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回妙传失败
	if !fileMetaResp.Suc {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}

	// 4. 上传过则将文件信息写入用户文件表，返回成功
	fmeta := dbClient.TableFileToFileMeta(fileMetaResp.Data.(orm.TableFile))
	fmeta.FileName = filename
	upRes, err := dbClient.OnUserFileUploadFinished(username, fmeta)
	if err != nil && upRes.Suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
	return
}
