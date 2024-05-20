package logic

import (
	"backend/model"
	"backend/mysql"

	"go.uber.org/zap"
)

func GetAllCommunity() (data []*model.Community, err error) {
	if data, err = mysql.GetAllCommunity(); err != nil {
		zap.L().Error("调用mysql.GetAllCommunity后出错", zap.Error(err))
	}
	return
}

func GetCommunityDetailByID(cid uint64) (data *model.CommunityDetail, err error) {
	data, err = mysql.GetCommunityDetailByID(cid)
	if err != nil {
		zap.L().Error("调用mysql.GetCommunityDetailByID后出错", zap.Error(err))
	}
	return
}
