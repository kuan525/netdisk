package mq

import "log"

var done chan bool

// StartConsume 接受消息
func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  // 自动应答
		false, // 非唯一的消费者
		false, // rabbitMQ只能设置为false
		false, // noWait， false表示回阻塞直到有消息过来
		nil)
	if err != nil {
		log.Fatal(err.Error(), "接受消息失败")
	}

	done = make(chan bool)

	go func() {
		// 循环读取channel的数据
		for d := range msgs {
			processErr := callback(d.Body)
			if processErr {
				// TODO: 将任务写入错误队列，待后续处理
			}
		}
	}()

	// 接受done的信号，没有信息过来则会一直阻塞，避免改函数推出
	<-done

	// 关闭通道
	channel.Close()
}
