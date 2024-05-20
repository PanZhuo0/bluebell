package mysql

import (
	"backend/model"

	"go.uber.org/zap"
)

func GetAllCommunity() (data []*model.Community, err error) {
	data = make([]*model.Community, 100)
	sqlStr := `select community_id,community_name from community`
	if err = db.Select(&data, sqlStr); err != nil {
		zap.L().Error("MySQL执行:`select community_id,community_name from community`出错", zap.Error(err))
	}
	return
}

func GetCommunityDetailByID(cid uint64) (data *model.CommunityDetail, err error) {
	data = new(model.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id =  ?`
	err = db.Get(data, sqlStr, cid)
	if err != nil {
		zap.L().Error("MySQL:执行`select community_id,community_name,introduction,create_time from community where community_id =  ?`时出错",
			zap.Error(err))
	}
	return
}
