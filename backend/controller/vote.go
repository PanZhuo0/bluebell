package controller

import (
	"backend/logic"
	"backend/model"
	"backend/redis"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VoteHandler(c *gin.Context) {
	vd := new(model.VoteData)
	err := c.ShouldBind(vd)
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
	err = logic.Vote(uid, vd)
	if err != nil {
		if errors.Is(err, redis.ErrorVoteRepeated) {
			zap.L().Error("用户重复投票", zap.Any("投票请求参数", vd))
			ResponseError(c, http.StatusBadRequest, CodeVoteRepeated)
			return
		}
		if errors.Is(err, redis.ErrorVoteTimeExpired) {
			zap.L().Error("投票时间已过", zap.Any("投票请求参数", vd))
			ResponseError(c, http.StatusBadRequest, CodeVoteTimeExpired)
			return
		}
		zap.L().Error("调用logic.Vote()后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccess(c)
}

func VoteNumHandler(c *gin.Context) {

}
