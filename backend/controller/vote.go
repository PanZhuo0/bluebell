package controller

import (
	"backend/logic"
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VoteHandler(c *gin.Context) {
	// 参数校验
	voteData := new(model.VoteData)
	err := c.ShouldBindJSON(voteData)
	if err != nil {
		zap.L().Error("请求参数有误", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	// 从context获取用户ID
	uid, err := GetUserID(c)
	if err != nil {
		zap.L().Error("用户用户ID失败", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeServeBusy)
		return
	}
	// 业务逻辑
	err = logic.Vote(uid, voteData)
	if err != nil {
		zap.L().Error("调用logic.Vote()后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccess(c)
}
