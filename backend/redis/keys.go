package redis

// redis中会记录每个帖子的投票数、帖子的创建时间、帖子的分数

const PrefixKey string = "bluebell:"

const (
	KeyZSetPostIDScore      = "post:score"
	KeyZSetPostIDCreateTime = "post:createTime"
	KeyZSetPostVotePrefix   = "post:vote:"
)

func getKey(k string) string {
	return PrefixKey + k
}
