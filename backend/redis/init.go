package redis

import (
	"backend/settings"
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var client *redis.Client

func Init(conf *settings.RedisConfig) {
	// 需要关闭虚拟机的防火墙(不知道应该是那个端口、或许是docker对应的那个、同时需要redis容器使用主机端口映射)
	client = redis.NewClient(&redis.Options{
		Addr: conf.Host + ":" + conf.Port,
		DB:   conf.DB,
	})
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Error("初始化Redis失败")
		panic(err)
	}
	zap.L().Info("初始化Redis成功", zap.String("Result", result))
}
