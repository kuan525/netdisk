package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	// Cli : redis client 客户端
	Cli       *redis.Client
	redisHost = "127.0.0.1:6379"
	redisPwd  = "redis_112525"
	// ctx 空上下文
	ctx = context.Background()
)

func init() {
	Cli := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPwd, // 如果有密码，请添加
		PoolSize: 50,       // v8底下自带连接池，这里设置容量为50
	})
	// 检查连接是否成功
	_, err := Cli.Ping(ctx).Result()
	if err != nil {
		log.Panic(err.Error(), "redis初始化失败")
	}
}
