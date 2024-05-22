package logic

import (
	"backend/model"
	"backend/mysql"
	"backend/redis"
	"backend/utils"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

var (
	ErrorIDsEmpty = errors.New("Redis中查询不到对应ID")
)

func CreatePost(uid uint64, p *model.Post) (err error) {
	// 生成PostID
	pid, err := utils.GenID()
	if err != nil {
		zap.L().Error("调用utils.GenID()后出错", zap.Error(err))
		return
	}
	p.PostID = pid

	// 生成帖子
	err = mysql.CreatePost(uid, p)
	if err != nil {
		zap.L().Error("调用mysql.CreatePost后出错", zap.Error(err))
	}
	// redis初始化帖子相关的数据
	err = redis.CreatePost(p)
	if err != nil {
		zap.L().Error("调用redis.CreatePost后出错", zap.Error(err), zap.Uint64("uid", uid), zap.Any("Post", p))
		return
	}
	return
}

func GetPostDetailByID(pid uint64) (data *model.APIPostDetail, err error) {
	data = new(model.APIPostDetail)
	var post *model.Post
	var user *model.User
	var communityDetail *model.CommunityDetail
	post, err = mysql.GetPostDetailByID(pid)
	// get post
	if err != nil {
		zap.L().Error("调用mysql.GetPostDetailByID()后出错", zap.Error(err), zap.Uint64("pid", pid))
		return
	}
	// get UserName by user_id
	user, err = mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("调用mysql.GetUserByID()后出错", zap.Uint64("post.AuthorID", post.AuthorID))
		return
	}
	// get CommunityDetail by community_id
	communityDetail, err = mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("调用mysql.GetCommunityDetailByID()后出错", zap.Uint64("post.CommunityID", post.CommunityID))
		return
	}
	data.Post = post
	data.AuthorName = user.UserName
	data.CommunityDetail = communityDetail
	return
}

func GetPostList(p *model.ParamPostList) (data []*model.APIPostDetail, err error) {
	// 1.获取redis中的ID
	ids, err := redis.GetPostListIDs(p)
	if err != nil {
		zap.L().Error("调用redis.GetPostListIDs()后出错", zap.Error(err), zap.Any("参数", p))
		return
	}
	// Redis中无对应ID数据
	if len(ids) == 0 {
		return nil, ErrorIDsEmpty
	}
	// 2.从MySQL中取出数据
	// 方法一:通过sqlx.In 批量获取post数据, 再根据每个POST的community_ID,author_id获取(这似乎有点画蛇添足,但是我还是实现了ids sqlx.In 批量查询)
	// posts, err := mysql.GetPostListByIDs(ids) //还需实现通过Post_id 查询对应社区、用户的信息,返回APIPostDetail

	// 方法二，直接for 循环调用logic.GetPostDetailByID()
	for i := 0; i < len(ids); i++ {
		id, err := strconv.ParseInt(ids[i], 10, 64)
		if err != nil {
			zap.L().Error("获取第index个帖子的ID时出错", zap.Error(err), zap.Int("Index", i))
			return nil, err
		}
		p, err := GetPostDetailByID(uint64(id))
		if err != nil {
			zap.L().Error("获取第index个帖子的详细信息时出错", zap.Error(err), zap.Int("Index", i))
			return nil, err
		}
		data = append(data, p)
	}
	return
}

func GetPostListByCommunity(p *model.ParamPostList) (posts []*model.APIPostDetail, err error) {
	// 从redis中获取对应IDs
	ids, err := redis.GetPostListIDsByCommunity(p)
	if err != nil {
		zap.L().Error("调用redis.GetPostListIDsByCommunity()后出错", zap.Error(err), zap.Any("参数", p))
		return
	}
	// 依照IDs去MySQL中查询数据
	if len(ids) == 0 { //如果获取到的IDs数组为空
		zap.L().Info("对应社区该页中无数据")
		return nil, ErrorIDsEmpty
	}
	// 循环获取(这样性能不如预处理快)
	for i := 0; i < len(ids); i++ {
		id, err := strconv.ParseInt(ids[i], 10, 64)
		if err != nil {
			zap.L().Error("解析第index个ID时出错", zap.Error(err), zap.Int("index", i))
			return nil, err
		}
		post, err := GetPostDetailByID(uint64(id))
		if err != nil {
			zap.L().Error("调用GetPostDetailByID()获取帖子详情时出错", zap.Error(err), zap.Int("index", i))
			return nil, err
		}
		posts = append(posts, post)
	}
	return
}
