package controller

import "github.com/gin-gonic/gin"

func ATestHandler(c *gin.Context) {
	c.String(200, "nihao")
}
