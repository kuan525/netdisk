package common

import (
	"github.com/urfave/cli/v2"
)

// CustomFlags 自定义命令行参数
var CustomFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "db-host",
		Value: "127.0.0.1",
		Usage: "database address",
	},
	&cli.StringFlag{
		Name:  "mq-host",
		Value: "127.0.0.1",
		Usage: "mq(rabbitmq) address",
	},
	&cli.StringFlag{
		Name:  "cache-host",
		Value: "127.0.0.1",
		Usage: "cache(redis) address",
	},
	&cli.StringFlag{
		Name:  "ceph-host",
		Value: "127.0.0.1",
		Usage: "ceph address",
	},
}
