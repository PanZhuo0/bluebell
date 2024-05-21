package redis

import (
	"backend/model"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func CreatePost(p *model.Post) (err error) {
	// 在bluebell:post:createTime中增加该Post对应的记录
	// 整个创建过程应该在一个事务中完成
	pipe := client.TxPipeline()
	pipe.ZAdd(context.Background(), getKey(KeyZSetPostIDCreateTime), redis.Z{
		Member: p.PostID,
		Score:  float64(time.Now().Unix()), //当前时间戳
	}).Result()
	// 在bluebell:post:score中增加该Post对应的记录
	pipe.ZAdd(context.Background(), getKey(KeyZSetPostIDScore), redis.Z{
		Member: p.PostID,
		Score:  float64(time.Now().Unix()),
	}).Result()
	// 在bluebell:community:comunityID增加对应的记录
	cid := strconv.Itoa(int(p.CommunityID)) //社区id
	pipe.SAdd(context.Background(), getKey(KeySetCommunityPrefix)+cid, p.PostID).Result()
	_, err = pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error("Redis:执行创建Post事务时出错", zap.Error(err))
		return
	}
	return
}

func GetPostListIDs(p *model.ParamPostList) (ids []string, err error) {
	// 获取redis中对应的index
	start := (p.Page - 1) * p.Size
	end := p.Page*p.Size - 1
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

// 根据参数从指定社区取出按分数/时间排序的某处的ID
func GetPostListIDsByCommunity(p *model.ParamPostListInSpecialCommunity) (ids []string, err error) {
	// 0.获取对应缓存key
	cachekey := p.Order + strconv.Itoa(int(p.CommunityID)) //eg:time[12312313--社区ID]
	communityKey := getKey(KeySetCommunityPrefix + strconv.Itoa(int(p.CommunityID)))
	var listKey string
	if p.Order == model.OrderScore {
		listKey = getKey(KeyZSetPostIDScore)
	}
	if p.Order == model.OrderTime {
		listKey = getKey(KeyZSetPostIDCreateTime)
	}
	// 1.查看缓存key是否过期
	if client.Exists(context.Background(), cachekey).Val() < 1 {
		// 不存在(过期)则重新计算(产生)
		pipe := client.Pipeline()
		pipe.ZInterStore(context.Background(), cachekey, &redis.ZStore{
			Keys:      []string{communityKey, listKey}, //获取两个key的交集,保留两者中分数较大者的数据
			Aggregate: "MAX",                           //interstore (community:[cid] , 所有帖子按时间/热度排序的key)
		})
		pipe.Expire(context.Background(), cachekey, 60*time.Second)
		_, err := pipe.Exec(context.Background())
		if err != nil {
			zap.L().Error("执行创建社区排行cacheKEY时失败", zap.Any("参数", p))
			return nil, err
		}
	}
	// 2.不过期直接取缓存key中对应数据即可,过期则更新缓存器再取数据
	start := (p.Page - 1) * p.Size
	end := p.Page*p.Size - 1
	ids = client.ZRevRange(context.Background(), cachekey, int64(start), int64(end)).Val()
	if len(ids) == 0 {
		zap.L().Info("未能从缓存Key中获取数据",
			zap.String("缓存key", cachekey),
			zap.String("对应的社区key", communityKey),
			zap.String("对应的列表key", listKey))
		return
	}
	return
}
