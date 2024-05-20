package redis

import (
	"backend/model"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func CreatePost(p *model.Post) (err error) {
	// 在bluebell:post:createTime中增加该Post对应的记录
	_, err = client.ZAdd(context.Background(), getKey(KeyZSetPostIDCreateTime), redis.Z{
		Member: p.PostID,
		Score:  float64(time.Now().Unix()), //当前时间戳
	}).Result()
	if err != nil {
		zap.L().Error("Redis:ZAdd()失败", zap.Error(err),
			zap.Any("Post Object", p),
			zap.String("key", KeyZSetPostIDCreateTime))
		return
	}
	// 在bluebell:post:score中增加该Post对应的记录
	_, err = client.ZAdd(context.Background(), getKey(KeyZSetPostIDScore), redis.Z{
		Member: p.PostID,
		Score:  float64(time.Now().Unix()),
	}).Result()
	if err != nil {
		zap.L().Error("Redis:ZAdd()失败", zap.Error(err),
			zap.Any("Post Object", p),
			zap.String("key", KeyZSetPostIDScore))
		return
	}
	return
}
