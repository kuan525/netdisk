package process

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/kuan525/netdisk/client/dbproxy"
	"github.com/kuan525/netdisk/mq"
	"github.com/kuan525/netdisk/store/cos"
	"log"
	"os"
)

// Transfer 处理文件转移
func Transfer(msg []byte) bool {
	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()

	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	_, err = cos.Client().Object.Put(context.Background(), pubData.DestLocation, bufio.NewReader(fin), nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 更新文件位置
	resp, err := dbClient.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !resp.Suc {
		log.Println("更新数据库异常，请检查: " + pubData.FileHash)
		return false
	}
	return true
}
