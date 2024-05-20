package logic

import (
	"backend/model"
	"backend/mysql"
	"backend/redis"

	"go.uber.org/zap"
)

func CreatePost(uid uint64, p *model.Post) (err error) {
	err = mysql.CreatePost(uid, p)
	if err != nil {
		zap.L().Error("调用mysql.CreatePost后出错", zap.Error(err))
	}
	err = redis.CreatePost(uid, p)
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
	var community *model.CommunityDetail
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
	// get CommunityName by community_id
	community, err = mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("调用mysql.GetCommunityDetailByID()后出错", zap.Uint64("post.CommunityID", post.CommunityID))
		return
	}
	data.Post = post
	data.AuthorName = user.UserName
	data.CommunityName = community.CommunityName
	return
}

func GetPostList() (data []*model.Post, err error) {
	// 具体需求还不确定,暂时搁置
	return
}
