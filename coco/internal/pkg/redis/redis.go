package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

//通用redis context
var Ctx_redis context.Context

//连接redis
func ConnectRedis(redisip string,redisport string,redispasswd string,redisdb int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisip+":"+redisport,
		Password: redispasswd,
		DB:       redisdb,
	})
	Ctx_redis = context.Background()
	return client
}

