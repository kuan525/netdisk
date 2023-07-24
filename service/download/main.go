package main

import (
	"download/route"
	"download/rpc"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	downloadProto "github.com/kuan525/netdisk/client/download/proto"
	cfg "github.com/kuan525/netdisk/config"
	"os"
)

func startRPCService() {
	logger := log.NewStdLogger(os.Stdout)
	log := log.NewHelper(logger)

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer(
		//grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)

	downloadProto.RegisterDownloadServiceServer(grpcSrv, new(rpc.Download))

	r := consul.New(consulClient)
	app := kratos.New(
		kratos.Name("go.kratos.service.download"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
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
