package mq

import (
	"github.com/kuan525/netdisk/config"
	"github.com/streadway/amqp"
	"log"
)

var (
	conn        *amqp.Connection
	channel     *amqp.Channel
	notifyClose chan *amqp.Error
)

// Init 初始化mq连接信息
func Init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMq功能
	if !config.AsyncTransferEnable {
		return
	}

	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}

	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}

// UpdateRabbitHost 更新mq host
func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}

func initChannel(rabbitHost string) bool {
	if channel != nil {
		return true
	}

	conn, err := amqp.Dial(rabbitHost)
	if err != nil {
		log.Println(err.Error(), "拨号失败")
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Println(err.Error(), "获取管道失败")
		return false
	}

	return true
}
