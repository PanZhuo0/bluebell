package main

import (
	"backend/logger"
	"backend/mysql"
	"backend/redis"
	"backend/router"
	"backend/settings"
	"strconv"
)

func main() {
	// 1.初始化配置文件
	settings.Init()
	// 2.初始化日志
	logger.Init(settings.Conf.LogConfig)
	// 3.初始化mysql
	mysql.Init(settings.Conf.MySQLConfig)
	// 4.初始化Redis
	redis.Init(settings.Conf.RedisConfig)
	// 5.初始化Router
	r := router.Setup()
	// 6.运行Router
	r.Run(":" + strconv.Itoa(settings.Conf.Port))
}
