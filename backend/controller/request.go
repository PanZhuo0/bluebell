package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ContextUserID = "USERID"

func GetUserID(c *gin.Context) (uid uint64, err error) {
	uidAny, ok := c.Get(ContextUserID)
	if !ok {
		zap.L().Error("上下文中没有USERID信息")
		ResponseError(c, http.StatusBadRequest, CodeContextWrong)
		return
	}
	// 类型断言与转换
	uidStr := uidAny.(string)
	uidInt64, err := strconv.ParseInt(uidStr, 10, 64) //str -> int64
	if err != nil {
		zap.L().Error("strconv.ParseInt()出错", zap.Error(err), zap.String("uidStr", uidStr))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	uid = uint64(uidInt64) //int64 -> uint64
	return
}
