package middleware

import (
	"backend/controller"
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ContextUserID = "USERID"

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		// the cases below declare token is invalid
		if authHeader == "" {
			controller.ResponseError(c, http.StatusBadRequest, controller.CodeNeedToken)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseErrorWithMsg(c, http.StatusBadRequest, controller.CodeInvalidToken, "Token格式不对")
			c.Abort()
			return
		}
		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			zap.L().Error("解析Token时出错", zap.String("Token:", parts[1]))
			controller.ResponseErrorWithMsg(c, http.StatusBadRequest, controller.CodeInvalidToken, err.Error())
			c.Abort()
			return
		}
		// Token有效,把Token中的ID设置到上下文中
		c.Set(ContextUserID, mc.UserID)
		c.Next()
	}
}
