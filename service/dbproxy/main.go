package main

import (
	dbRpc "dbproxy/rpc"
	"github.com/kuan525/netdisk/common"
	"github.com/kuan525/netdisk/config"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"log"

	dbConn "github.com/kuan525/netdisk/dbclient/conn"
	dbproto "github.com/kuan525/netdisk/dbclient/proto"
	"time"
)

func startRpcService() {
	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"), // 在注册中心中等服务名称
		micro.RegisterTTL(time.Second*10),      // 声明超时时间，避免consul不主动删除掉已经失去心跳等服务节点
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			//检查是否指定dbhost
			dbhost := c.String("dbhost")
			if len(dbhost) > 0 {
				log.Println("custom db address: " + dbhost)
				config.UpdateDBHost(dbhost)
			}
		}),
	)

	// 初始化db connection
	dbConn.InitDBConn()

	dbproto.RegisterDBProxyServiceHandler(service.Server(), new(dbRpc.DBProxy))
	if err := service.Run(); err != nil {
		log.Println(err.Error())
	}
}

func main() {
	startRpcService()
}
