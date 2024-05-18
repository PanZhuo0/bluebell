package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func Init() {
	// 需要关闭虚拟机的防火墙(不知道应该是那个端口、或许是docker对应的那个、同时需要redis容器使用主机端口映射)
	client = redis.NewClient(&redis.Options{
		Addr: "192.168.18.133:6379", //虚拟机的IP
		DB:   0,
	})
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Redis 初始化成功！", result)
}
