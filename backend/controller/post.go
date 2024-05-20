package controller

import (
	"backend/logic"
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// 1.1参数校验
	p := new(model.Post)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("请求参数有误", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	// 1.2从context获取用户id
	uid, err := GetUserID(c)
	if err != nil {
		return
	}
	// 业务逻辑
	err = logic.CreatePost(uid, p)
	if err != nil {
		zap.L().Error("调用logic.CreatePost()出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccess(c)
}
