package logic

import (
	"backend/model"
	"backend/redis"

	"go.uber.org/zap"
)

func Vote(uid uint64, vd *model.VoteData) (err error) {
	err = redis.Vote(uid, vd)
	if err != nil {
		zap.L().Error("调用redis.Vote()后出错", zap.Error(err))
		return
	}
	return
}
