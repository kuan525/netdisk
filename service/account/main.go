package main

import (
	"account/handler"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	userProto "github.com/kuan525/netdisk/client/account/proto"
	"os"
)

func main() {
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

	userProto.RegisterUserServiceServer(grpcSrv, new(handler.User))

	r := consul.New(consulClient)
	app := kratos.New(
		kratos.Name("go.kratos.service.user"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
