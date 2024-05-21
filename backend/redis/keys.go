package redis

// redis中会记录每个帖子的投票数、帖子的创建时间、帖子的分数

const PrefixKey string = "bluebell:"

const (
	KeyZSetPostIDScore      = "post:score"
	KeyZSetPostIDCreateTime = "post:createTime"
	KeyZSetPostVotePrefix   = "post:vote:"
	KeySetCommunityPrefix   = "community:" //保存每个分区下,帖子的ID ,community:go([community_id]) 这个key 下有各种帖子
	// 每次创建帖子的时候,都把对应的数据增加到community:[community_id] set中
	// 可以通过Set与ZSet的联合存储 zinterstore 获得两者这并有的数据
	// 需要利用缓存key 减少zinterstore的访问次数
)

func getKey(k string) string {
	return PrefixKey + k
}
