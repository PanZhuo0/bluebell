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
const VoteDuration = 3600 * 24 * 7

var (
	ErrorVoteTimeExpired      = errors.New("投票时间已过")
	ErrorPostIDRecordNotExist = errors.New("对应PostID记录不存在")
	ErrorVoteRepeated         = errors.New("重复投票")
)

func Vote(uid uint64, vd *model.VoteData) (err error) {
	// 0.解析数据
	pidStr := fmt.Sprintf("%d", vd.PostID)
	uidStr := fmt.Sprintf("%d", uid)
	// 1.是否超过可投票时间?
	postCreateTime := client.ZScore(context.Background(), getKey(KeyZSetPostIDCreateTime), pidStr).Val()
	if time.Now().Unix()-int64(postCreateTime) > VoteDuration {
		//如果超时
		return ErrorVoteTimeExpired
	}
	// 2.是否重复投票?
	ov := client.ZScore(context.Background(), getKey(KeyZSetPostVotePrefix)+pidStr, uidStr).Val()
	diff := vd.Direction - ov
	if diff == 0 {
		return ErrorVoteRepeated
	}
	// 投票(一个事务)
	// 3.1修改帖子的分数
	pipe := client.TxPipeline()
	pipe.ZIncrBy(context.Background(), getKey(KeyZSetPostIDScore), diff*scorePerVote, pidStr) //post_scoreZSET(member:pid score:增加)
	// 3.2修改帖子对应用户的投票记录
	pipe.ZAdd(context.Background(), getKey(KeyZSetPostVotePrefix)+pidStr, redis.Z{ // vote:post_id  (Member:uid score:赞同1/反对-1)
		Member: uidStr,
		Score:  vd.Direction,
	})
	// 3.3执行事务
	_, err = pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error("Redis:执行投票事务时出错", zap.Error(err))
		return
	}
	return
}
