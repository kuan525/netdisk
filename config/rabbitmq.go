package config

const (
	// AsyncTransferEnable 异步转移
	AsyncTransferEnable  = false
	TranExchangeName     = "uploadserver.trans"
	TransCOSQueueName    = "uploadserver.trans.cos"
	TransCOSErrQueueName = "uploadserver.trans.cos.err"
	TransCOSRoutingKey   = "cos"
)

var (
	RabbitURL = "amqp://guest:guest@127.0.0.1:5672/"
)
