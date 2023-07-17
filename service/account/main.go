package main

import (
	"account/handler"
	"github.com/kuan525/netdisk/common"
	dbcli "github.com/kuan525/netdisk/dbclient"
	accProto "github.com/kuan525/netdisk/proto/account"
	"github.com/micro/go-micro"
	"log"
	"time"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	// 初始化service，解析命令行参数等
	service.Init()

	// 初始化dbclient
	dbcli.Init(service)

	accProto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err.Error())
	}
}
