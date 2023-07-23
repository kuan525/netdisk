package main

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	dbConn "github.com/kuan525/netdisk/service/dbproxy/conn"
	dbRpc "github.com/kuan525/netdisk/service/dbproxy/rpc"
	"os"

	dbproto "github.com/kuan525/netdisk/client/dbproxy/proto"
)

func startRpcService() {
	logger := log.NewStdLogger(os.Stdout)
	log := log.NewHelper(logger)

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer(
		//grpc.Address(":9000"), // 默认随机
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)
	dbproto.RegisterDBProxyServiceServer(grpcSrv, new(dbRpc.DBProxy))

	r := consul.New(consulClient)
	app := kratos.New(
		kratos.Name("go.micro.service.dbproxy"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)

	// 初始化db connection
	dbConn.InitDBConn()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startRpcService()
}
