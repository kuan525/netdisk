package main

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	upProto "github.com/kuan525/netdisk/client/upload/proto"
	cfg "github.com/kuan525/netdisk/config"
	"github.com/kuan525/netdisk/mq"
	"os"
	"upload/route"
	upRpc "upload/rpc"
)

func startRPCService() {
	logger := log.NewStdLogger(os.Stdout)
	log := log.NewHelper(logger)

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer(
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)
	upProto.RegisterUploadServiceServer(grpcSrv, new(upRpc.Upload))

	r := consul.New(consulClient)
	app := kratos.New(
		kratos.Name("go.kratos.service.upload"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)

	// 初始化mq client
	mq.Init()

	if err := app.Run(); err != nil {
		log.Fatal(err)
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
	startRPCService()
}
