package main

import (
	"backend/logger"
	"backend/redis"
	"backend/settings"
)

func main() {
	// 1.初始化配置文件
	settings.Init()
	// 2.初始化日志
	logger.Init(settings.Conf.LogConfig)
	// 3.初始化mysql
	// mysql.init(settings.Conf.MySQLConfig)
	// 4.初始化Redis
	redis.Init()
	// 5.初始化Router
	// 6.运行Router
}
