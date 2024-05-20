package logic

import (
	"backend/model"
	"backend/mysql"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrorUserExist = errors.New("用户已存在")
)

func SignUp(p *model.ParamSignUp) (err error) {
	// 1.user exists?
	exist, err := mysql.UserCheckExist(p)
	if err != nil {
		zap.L().Error("调用mysql.UserCheckExist()时出错", zap.Error(err))
		return
	}
	if exist {
		zap.L().Error("用户已存在", zap.String("用户名", p.UserName))
		return ErrorUserExist
	}
	// 2.user register user
	if err = mysql.UserRegister(p); err != nil {
		zap.L().Error("调用 mysql.UserRegister() 时出错", zap.Error(err))
	}
	return
}

func Login(u *model.User) (err error) {
	// 1.验证用户账号密码
	if err = mysql.CheckUserPassword(u); err != nil {
		zap.L().Error("调用mysql.CheckUserPassword后出错", zap.Error(err))
	}
	return
}
