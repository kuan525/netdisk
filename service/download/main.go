package main

import (
	"fmt"
	"github.com/kuan525/netdisk/common"
	dbcli "github.com/kuan525/netdisk/dbclient"
	dlProto "github.com/kuan525/netdisk/proto/download"
	cfg "github.com/kuan525/netdisk/service/download/config"
	"github.com/kuan525/netdisk/service/download/route"
	dlRpc "github.com/kuan525/netdisk/service/download/rpc"
	"github.com/micro/go-micro"
	"time"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init()

	// 初始化dbclient
	dbcli.Init(service)

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main() {
	// api服务
	go startAPIService()

	// rpc服务
	startRPCService()
}
