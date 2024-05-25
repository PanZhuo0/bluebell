package main

import (
	"backend/logger"
	"backend/mysql"
	"backend/redis"
	"backend/router"
	"backend/settings"
	"backend/utils"
	"strconv"
)

const MachineID = 1000

// @title bulebell项目
// @Version 1.0
// @description 这里是描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.utl http://www.swagger.io/support
// @contrct.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/license/LICENSE-2.0.html

// @host 这里写host主机
// @BasePath 这里写base path 根地址
func main() {
	// pre:初始化工具
	utils.Init(MachineID)
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
