package cache

import (
	"context"
	"goproxy/pkg/config"
	"log"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

//初始化连接
func init(){	
	Redis = redis.NewClient(&redis.Options{Addr: config.GetString("cache.redis")})
    pong, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error Connecting Redis :%s", err)
	}
	log.Printf("Dial to redis->%s>>%s", config.GetString("cache.redis"), pong)	
}
