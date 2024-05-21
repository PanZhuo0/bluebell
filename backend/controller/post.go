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

// 根据PostID 获取Post的详细信息(查询用户表、社区表)
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

// Addition:需要获取每个帖子中的赞同数
// Addition:增加一个定时任务、如果post的时间超过一小时、停止投票（把post从redis中的对应zset中删除、分数表、time/score表）
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

// 获取某个社区内的帖子列表(允许按时间/热度排序,需要分页)
func GetPostListByCommunityHandler(c *gin.Context) {
	// 1.参数获取与处理
	p := &model.ParamPostListInSpecialCommunity{
		CommunityID: 0,
		Page:        1,
		Size:        10,
		Order:       model.OrderTime,
	}
	err := c.ShouldBindQuery(p)
	if err != nil || (p.Order != model.OrderScore && p.Order != model.OrderTime) || p.Page <= 0 {
		zap.L().Error("请求参数有误", zap.Any("请求参数", p), zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	// 2.业务逻辑
	data, err := logic.PostListByCommunity(p)
	if err != nil {
		if errors.Is(err, logic.ErrorIDsEmpty) {
			zap.L().Info("从Redis中并未获取到IDs", zap.Error(err), zap.Any("参数", p))
			ResponseError(c, http.StatusAccepted, CodePageInvalid)
			return
		}
		zap.L().Error("调用logic.PostListByCommunity()后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	// 响应
	ResponseSuccessWithData(c, data)
}
