package router

import (
	"backend/controller"
	"backend/controller/middleware"
	"time"

	// _ "bluebell/docs"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.RateLimitMiddleware(time.Second*2, 1)) //限流，令牌桶，2秒产生1个令牌，最多1个
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.Use(middleware.AuthMiddleWare())
	{
		v1.GET("/communities", controller.CommunityHandler)
		v1.GET("/communityDetail/:id", controller.CommunityDetailHandler)

		v1.GET("/post", controller.PostListHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.PostDetailHandler)

		v1.POST("/vote", controller.VoteHandler)
	}
	return r
}
