package main

import (
	"backend/redis"
	"backend/settings"
)

func main() {
	// 1.初始化配置文件
	settings.Init()
	// 2.初始化日志
	// 3.初始化mysql
	// 4.初始化Redis
	redis.Init()
	// 5.初始化Router
	// 6.运行Router
}
