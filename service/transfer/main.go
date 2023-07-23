package main

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	"github.com/kuan525/netdisk/config"
	"github.com/kuan525/netdisk/mq"
	"os"
	"transfer/process"
)

func startRPCService() {
	logger := log.NewStdLogger(os.Stdout)
	log := log.NewHelper(logger)

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	r := consul.New(consulClient)
	app := kratos.New(
		kratos.Name("go.micro.service.transfer"),
		kratos.Server(),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func startTransferService() {
	if !config.AsyncTransferEnable {
		log.Info("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Info("文件转移服务启用中，开始监听转移任务队列")

	// 初始化mq client
	mq.Init()

	mq.StartConsume(config.TransCOSQueueName, "transfer_cos", process.Transfer)
}

func main() {
	// 文件转移服务
	go startTransferService()

	// rpc服务
	startRPCService()

}
