package main

import (
	"apigw/route"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	"os"
)

// startRPCService 其实没啥用，在consul中占个名儿
func startRPCService() {
	logger := log.NewStdLogger(os.Stdout)
	log := log.NewHelper(logger)

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	r := consul.New(consulClient)
	app := kratos.New(
		kratos.Name("go.micro.service.apigw"),
		kratos.Server(),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	go startRPCService()

	r := route.Router()
	r.Run("localhost:8080")
}
