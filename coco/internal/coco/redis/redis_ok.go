package coco

import (
	"coco/internal/coco"
	myconfig "coco/internal/coco/config"
	mydb "coco/internal/coco/mongo"
	myredis "coco/internal/pkg/redis"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
)

var RedisClient *redis.Client

//redis初始化
func InitRedis() {
	RedisClient = myredis.ConnectRedis(myconfig.ConfigStore.RedisIp, myconfig.ConfigStore.RedisPort, myconfig.ConfigStore.RedisPasswd, myconfig.ConfigStore.RedisDB)
	if RedisClient != nil {
		fmt.Println("Redis Success!",RedisClient)
		return
	}
	fmt.Println("Failed to Connect Redis!")
	os.Exit(1)
}

//向redis中增添数据
func AddUserToRedis(key string, v1, v2 interface{}) error {
	_, err := RedisClient.HSet(myredis.Ctx_redis, key, v1, v2).Result()
	if err != nil {
		return coco.AddRedisMapError
	}
	return nil
}

//到redis中查找数据
func FindUserInRedis(key string, v1 string, v2 string) (bool, error) {
	b, err := RedisClient.HGet(myredis.Ctx_redis, key, v1).Result()
	if err != nil {
		return false, nil
	}
	//检查用户是否被禁止登陆
	if CheckLoginRedis(v1) {
		return false, coco.ForbiddenLoginError
	}
	//检查用户是否已经登陆
	err = CheckUserHasLogin(iris.Map{"username": v1})
	if err != nil {
		return false, coco.UserHasLoginError
	}
	if b != v2 {
		return false, coco.PasswordError
	}
	return true, nil
}

//密码错误后检测，3次封锁账号5分钟
func CheckPasswdRedis(key string, v1 string) {
	t, err := RedisClient.HGet(myredis.Ctx_redis, key, v1).Result()
	if err != nil {
		AddUserToRedis(key, v1, "1")
		return
	}
	if t == "1" {
		RedisClient.HIncrBy(myredis.Ctx_redis, key, v1, 1).Result()
		return
	}
	RedisClient.HDel(myredis.Ctx_redis, key, v1).Result()
	RedisClient.Set(myredis.Ctx_redis, v1, "1", time.Minute*time.Duration(myconfig.ConfigStore.UserPasswdErrLock))
}

//登陆时检测账号是否被禁止
func CheckLoginRedis(key string) bool {
	_, err := RedisClient.Get(myredis.Ctx_redis, key).Result()
	if err == nil {
		return true
	}
	return false
}

//检测用户是否已经登陆
func CheckUserHasLogin(data iris.Map) error {
	u := coco.WhiteStore{}
	err := mydb.MongoDB.Find(u.CollectName(), iris.Map{"username": data["username"]}, &u)
	if err == nil {
		return coco.UserHasLoginError
	}
	return nil
}

//登陆成功后清除forbidden中的无效记录
func CleanForbidden(key string, v1 string) {
	b, _ := RedisClient.HExists(myredis.Ctx_redis, key, v1).Result()
	if b {
		RedisClient.HDel(myredis.Ctx_redis, key, v1).Result()
	}
}
