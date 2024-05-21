package redis

import (
	"backend/model"
	"context"
	"fmt"
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

func GetPostListIDs(p *model.ParamPostList) (ids []string, err error) {
	// 获取redis中对应的index
	start := (p.Page - 1) * p.Size
	end := p.Page*p.Size - 1
	fmt.Println(start)
	fmt.Println(end)
	if p.Order == model.OrderTime {
		ids = client.ZRevRange(context.Background(), getKey(KeyZSetPostIDCreateTime), int64(start), int64(end)).Val()
		return
	}
	if p.Order == model.OrderScore {
		ids = client.ZRevRange(context.Background(), getKey(KeyZSetPostIDScore), int64(start), int64(end)).Val()
		return
	}
	zap.L().Error("请求参数有误")
	return
}
