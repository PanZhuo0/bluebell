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

// CreatePostHandler 创建帖子接口
// 简介
// @Summary 创建帖子接口
// 描述信息
// @Description 可以用于创建帖子
// 分组
// @Tags 帖子相关接口
// 接收的内容类型
// @Accpet application/json
// 返回的内容类型
// @Produce application/json
// 需要的参数,true:必填
// @Param Authorization header string true "Bearer 用户令牌"
// 需要的参数,query:放在Query
// @Param object query model.Post false "查询参数"
// 安全性:API需要认证
// @Security ApiKeyAuth
// 成功返回ResData对象类型的数据
// @Success 200 {object} ResData
// 接收的请求类型和请求路径
// @Router /post [post]
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
func PostListHandler(c *gin.Context) { // 整合:按社区查询、全部查询
	// 1.参数获取与处理
	p := &model.ParamPostList{
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
	var data []*model.APIPostDetail
	if p.CommunityID == 0 {
		// (communityID==0)未指定社区,将查询全部帖子
		data, err = logic.GetPostList(p)
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
	} else {
		// 将查询指定社区
		data, err = logic.GetPostListByCommunity(p)
		if err != nil {
			if errors.Is(err, logic.ErrorIDsEmpty) {
				zap.L().Info("从Redis中并未获取到IDs", zap.Error(err), zap.Any("参数", p))
				ResponseError(c, http.StatusAccepted, CodePageInvalid)
				return
			}
			zap.L().Error("调用logic.GetPostListByCommunity()后出错", zap.Error(err))
			ResponseError(c, http.StatusAccepted, CodeServeBusy)
			return
		}
		// 响应
		ResponseSuccessWithData(c, data)
	}
}
