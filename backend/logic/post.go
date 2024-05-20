package logic

import (
	"backend/model"
	"backend/mysql"

	"go.uber.org/zap"
)

func CreatePost(uid uint64, p *model.Post) (err error) {
	err = mysql.CreatePost(uid, p)
	if err != nil {
		zap.L().Error("调用mysql.CreatePost后出错", zap.Error(err))
	}
	return
}
