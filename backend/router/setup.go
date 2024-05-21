package router

import (
	"backend/controller"
	"backend/controller/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
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
		v1.GET("/vote/:id", controller.VoteNumHandler)
	}
	return r
}
