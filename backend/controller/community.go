package controller

import (
	"backend/logic"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CommunityHandler(c *gin.Context) {
	// 业务逻辑
	data, err := logic.GetAllCommunity()
	if err != nil {
		zap.L().Error("调用logic.GetAllCommunity后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	// 3.响应数据
	ResponseSuccessWithData(c, data)
}

// 从URI获取参数ID
func CommunityDetailHandler(c *gin.Context) {
	// 参数获取
	cid := c.Param("id")
	ccid, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		zap.L().Error("请求参数有误", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	communityID := uint64(ccid)
	// 业务逻辑
	data, err := logic.GetCommunityDetailByID(communityID)
	if err != nil {
		zap.L().Error("调用getCommunityDetailByID后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccessWithData(c, data)
}
