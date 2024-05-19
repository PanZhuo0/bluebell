package router

import (
	"backend/controller"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.GET("/test", controller.ATestHandler)
	return r
}
