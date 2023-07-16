package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	rPool "github.com/kuan525/netdisk/cache/redis"
	"github.com/kuan525/netdisk/config"
	dbcli "github.com/kuan525/netdisk/dbclient"
	"github.com/kuan525/netdisk/util"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func init() {
	os.MkdirAll(config.TempLocalRootDir, 0744)
}

// InitialMultipartUploadHandler 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "params invalid",
			})
		return
	}

	// 2. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 3. 将初始化信息写入到redis缓存
	rPool.Cli.HSet(context.Background(), "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rPool.Cli.HSet(context.Background(), "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rPool.Cli.HSet(context.Background(), "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 4. 将响应初始化数据返回到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": upInfo,
		})
	return
}

// UploadPartHandler 上传文件分块
func UploadPartHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	// username := c.Request.FormValue("username")
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 2. 获取文件句柄，用于存储分块的内容
	fpath := config.TempPartRootDir + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.JSON(http.StatusOK,
			gin.H{
				"code": 0,
				"msg":  "upload part failed",
				"data": nil,
			})
		return
	}

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 3. 更新redis缓存状态
	rPool.Cli.HSet(context.Background(), "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 4. 返回处理结果到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})
	return
}

func CompleteUploadHandler(c *gin.Context) {
	// 1. 解析请求参数
	upid := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize := c.Request.FormValue("filesize")
	filename := c.Request.FormValue("filename")

	// 2. 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := rPool.Cli.HGetAll(context.Background(), "MP_"+upid).Result()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "服务错误",
				"data": nil,
			})
		return
	}

	totalCount, chunkCount := 0, 0
	for k, v := range data {
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "分块不完整",
				"data": nil,
			})
		return
	}

	// 3. TODO 合并分块，可以将ceph当临时存储，合并时将文件写入ceph
	// 也可以不用在本地进行合并，转移的时候将分块append到ceph/cos即可
	srcPath := config.TempPartRootDir + upid + "/"
	destPath := config.TempLocalRootDir + filehash
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
	mergeRes, err := util.ExecLinuxShell(cmd)
	if err != nil {
		log.Println(err.Error())
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "合并失败",
				"data": nil,
			})
		return
	}
	log.Println(mergeRes)

	// 4. 更新唯一文件表以及用户文件表
	fsize, _ := strconv.Atoi(filesize)

	fmeta := dbcli.FileMeta{
		FileSha1: filehash,
		FileName: filename,
		FileSize: int64(fsize),
		Location: destPath,
	}
	_, ferr := dbcli.OnFileUploadFinished(fmeta)
	_, uferr := dbcli.OnUserFileUploadFinished(username, fmeta)
	if ferr != nil || uferr != nil {
		log.Println(err.Error())
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "数据库更新失败",
				"data": nil,
			})
		return
	}

	// 5. 响应处理结果
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": nil,
		})
	return
}
