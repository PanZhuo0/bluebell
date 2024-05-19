package router

import (
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Group("/api/v1")
	return r
}
