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

// 获取数据？
// 将帖子对应的作者和社区查询出来，一并返回
// sqlx.In 如何使用?
// db.Rebind()是用来干嘛的？
// 需要模仿PostDetailByID 写个getpostlistby ids
// Addition:需要获取每个帖子中的赞同数
// Addition:增加一个定时任务、如果post的时间超过一小时、停止投票（把post从redis中的对应zset中删除、分数表、time/score表）
// Addition:解决投票赞成、现在在还投赞成的问题、需要增加个反馈给用户 if newvalue=oldvalue -->报个不能重复投票的错误
func PostListHandler(c *gin.Context) {
	// 参数检验
	p := &model.ParamPostList{
		Page:  1,
		Size:  10,
		Order: model.OrderTime,
	}
	err := c.ShouldBindQuery(p)
	if (p.Order != model.OrderScore && p.Order != model.OrderTime) || err != nil || p.Page <= 0 {
		zap.L().Error("请求参数有误", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	// 获取数据
	data, err := logic.GetPostList(p)
	if err != nil {
		zap.L().Error("调用logic.GetPostList()后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodePageInvalid)
		return
	}
	// 响应
	ResponseSuccessWithData(c, data)
}

// PostListInCommunityByTimes/Scores？
// 根据社区，在社区内的按time/score 排序(总体逻辑上来说是类似的，只是增加一个community_id的限制)
// 如何实现，使用zset 交集inter zset:
// 维护一个commmunity_list
