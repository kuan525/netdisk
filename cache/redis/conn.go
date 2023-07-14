package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	Cli       *redis.Client
	redisHost = "127.0.0.1:6379"
	redisPwd  = "redis_112525"
	ctx       = context.Background()
)

func init() {
	Cli := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPwd, // 如果有密码，请添加
		PoolSize: 50,
	})
	_, err := Cli.Ping(ctx).Result()
	if err != nil {
		log.Panic(err.Error(), "redis初始化失败")
	}
}
