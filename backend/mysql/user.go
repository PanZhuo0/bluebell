package mysql

import (
	"backend/model"
	"backend/utils"
	"errors"

	"go.uber.org/zap"
)

var ErrorPasswordWrong = errors.New("密码不正确")

func UserCheckExist(p *model.ParamSignUp) (exist bool, err error) {
	count := 0
	sqlStr := `select count(*) from user where username=?`
	err = db.Get(&count, sqlStr, p.UserName)
	if err != nil {
		zap.L().Error("执行 `select count(*) from user where username=?` 时出错", zap.Error(err))
		return
	}
	return count > 0, nil
}

func UserRegister(p *model.ParamSignUp) (err error) {
	var id uint64
	var password string
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	//产生用户ID
	if id, err = utils.GenID(); err != nil {
		zap.L().Error("调用GenID()后出错", zap.Error(err))
		return err
	}
	//用户密码加密
	password = utils.Encrypt(p.Password)
	if _, err = db.Exec(sqlStr, id, p.UserName, password); err != nil {
		zap.L().Error("MySQL:执行`insert into user(user_id,username,password) values(?,?,?)`时出错",
			zap.Error(err))
	}
	return
}

func CheckUserPassword(u *model.User) (err error) {
	var oPassword string
	// var inMySQLPassword
	sqlStr := `select user_id,username,password from user where username=?`
	oPassword = utils.Encrypt(u.Password)
	if err = db.Get(u, sqlStr, u.UserName); err != nil {
		zap.L().Error("MySQL:`select password from user where username=?`出错", zap.Error(err))
		return
	}
	if oPassword != u.Password { //at this moment,u.Password = inMySQLPassword
		return ErrorPasswordWrong
	}
	return
}
