package redis

import (
	"backend/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const scorePerVote = 432 // 86400/200=432
const Duration = 3600 * 24 * 7

var (
	ErrorVoteTimeExpired      = errors.New("投票时间已过")
	ErrorPostIDRecordNotExist = errors.New("对应PostID记录不存在")
)

func Vote(uid uint64, vd *model.VoteData) (err error) {
	// 是否超过可投票时间?
	uidStr := fmt.Sprintf("%d", uid)
	pidStr := fmt.Sprintf("%d", vd.PostID)
	key := getKey(KeyZSetPostVotePrefix) + pidStr
	postCreatedTime := client.ZScore(context.Background(), getKey(KeyZSetPostIDCreateTime), pidStr).Val()
	if postCreatedTime == 0 {
		zap.L().Error("对应PostID数据不存在,PostID正确吗?", zap.Any("VoteData", vd))
		return ErrorPostIDRecordNotExist
	}
	if float64(time.Now().Unix())-Duration > postCreatedTime {
		return ErrorVoteTimeExpired
	}
	// 需要建立事务完成下面的一系列任务
	pipe := client.Pipeline()
	// 修改对应帖子分数
	ov := pipe.ZScore(context.Background(), key, uidStr).Val()                                //原本的投票情况
	diff := vd.Direction - ov                                                                 //前后投票状态差值 (0:未投票,1:赞同,-1:反对)
	pipe.ZIncrBy(context.Background(), getKey(KeyZSetPostIDScore), diff*scorePerVote, pidStr) //修改帖子的积分

	// 修改用户投票状态
	_, err = pipe.ZAdd(context.Background(), key, redis.Z{
		Member: uid,
		Score:  vd.Direction,
	}).Result()
	if err != nil {
		zap.L().Error("Redis:ZAdd()出错", zap.Any("key", key), zap.Error(err))
		return
	}
	if _, err = pipe.Exec(context.Background()); err != nil {
		zap.L().Error("Redis:投票过程中出错", zap.Error(err))
		return
	}
	return
}
