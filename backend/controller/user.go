package controller

import (
	"backend/logic"
	"backend/model"
	"backend/mysql"
	"backend/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	// 1.获取参数,参数校验
	p := new(model.ParamSignUp)
	err := c.ShouldBindJSON(p)
	if err != nil {
		zap.L().Error("ParamSignUp绑定出错", zap.Error(err))
		ResponseError(c, http.StatusBadRequest, CodeInvalidParams)
		return
	}
	// 2.逻辑处理
	err = logic.SignUp(p)
	if err != nil {
		zap.L().Error("调用logic.SignUp()后出错", zap.Error(err))
		// 错误1.用户已存在,创建用户失败,查询用户时出错,
		if errors.Is(err, logic.ErrorUserExist) {
			ResponseError(c, http.StatusAccepted, CodeUserExist)
			return
		}
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	// 3.响应
	zap.L().Info("创建用户成功", zap.Any("user", p))
	ResponseSuccess(c)

}

func LoginHandler(c *gin.Context) {
	// 1.参数校验
	u := new(model.User)
	err := c.ShouldBind(u)
	if err != nil {
		zap.L().Error("请求参数有误", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeInvalidParams)
		return
	}
	// 2.业务逻辑
	if err = logic.Login(u); err != nil {
		zap.L().Error("调用logic.UserLogin后出错", zap.Error(err))
		if err == mysql.ErrorPasswordWrong {
			ResponseError(c, http.StatusAccepted, CodeUserPasswordWrong)
			return
		}
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	// 3.响应、这里需要返回TOKEN
	aToken, rToken, err := utils.GenToken(fmt.Sprintf("%d", u.UserID))
	if err != nil {
		zap.L().Error("调用utils.GenToken后出错", zap.Error(err))
		ResponseError(c, http.StatusAccepted, CodeServeBusy)
		return
	}
	ResponseSuccessWithData(c, gin.H{
		"accessToken":  aToken,
		"refreshToken": rToken,
		"userID":       u.UserID,
		"username":     u.UserName,
	})
}
