package transfer

import (
	"fmt"
	"github.com/kuan525/netdisk/common"
	"github.com/kuan525/netdisk/config"
	dbcli "github.com/kuan525/netdisk/dbclient"
	"github.com/kuan525/netdisk/mq"
	"github.com/kuan525/netdisk/service/transfer/process"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"log"
	"time"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro,service.transfer"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
		}),
	)

	// 初始化dbclient
	dbcli.Init(service)

	if err := service.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

func startTranserService() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启用中，开始监听转移任务队列")

	// 初始化mq client
	mq.Init()

	mq.StartConsume(config.TransCOSQueueName, "transfer_cos", process.Transfer)
}

func main() {
	// 文件转移服务
	go startTranserService()

	// rpc服务
	startRPCService()

}