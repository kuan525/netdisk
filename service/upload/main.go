package main

import (
	"fmt"
	"github.com/kuan525/netdisk/common"
	cfg "github.com/kuan525/netdisk/config"
	dbcli "github.com/kuan525/netdisk/dbclient"
	"github.com/kuan525/netdisk/mq"
	upProto "github.com/kuan525/netdisk/proto/upload"
	upRpc "github.com/kuan525/netdisk/service/upload/rpc"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	"log"
	"os"
	"time"
	"upload/route"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.upload"), // 服务名称
		micro.RegisterTTL(time.Second*10),     // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
		}))

	// 初始化dbproxy client
	dbcli.Init(service)
	mq.Init()

	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.UploadServiceHost)
}

func main() {
	os.MkdirAll(cfg.TempLocalRootDir, 0744)
	os.MkdirAll(cfg.TempPartRootDir, 0744)

	// api服务
	go startAPIService()

	// rpc 服务
	startAPIService()
}
