package controller

import (
	"backend/logic"
	"backend/model"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// parma check
	p := new(model.Post)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("请求参数有误", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	// get userid from Context
	uid, err := GetUserID(c)
	if err != nil {
		return
	}
	// logic
	err = logic.CreatePost(uid, p)
	if err != nil {
		zap.L().Error("调用logic.CreatePost()出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccess(c)
}

func PostDetailHandler(c *gin.Context) {
	// parma check and process
	pidStr := c.Param("id")
	pidInt64, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt()出错", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	pid := uint64(pidInt64)
	// logic
	data, err := logic.GetPostDetailByID(pid)
	if err != nil {
		zap.L().Error("调用logic.GetPostDetailByID()后出错",
			zap.Error(err),
			zap.Uint64("pid", pid))
		if errors.Is(err, sql.ErrNoRows) { //if not row be founded
			ResponseError(c, http.StatusAccepted, CodeDataNotExist)
			return
		}
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccessWithData(c, data)
}

func PostListHandler(c *gin.Context) {
	data, err := logic.GetPostList()
	if err != nil {
		zap.L().Error("调用logic.GetPostList()后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccessWithData(c, data)
}
