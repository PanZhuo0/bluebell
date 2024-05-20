package mysql

import (
	"backend/model"

	"go.uber.org/zap"
)

func CreatePost(uid uint64, p *model.Post) (err error) {
	// 1.产生POSTID

	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	if _, err = db.Exec(sqlStr, p.PostID, p.Title, p.Content, uid, p.CommunityID); err != nil {
		zap.L().Error("MySQL:执行`insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`时出错", zap.Error(err))
	}
	return
}

func GetPostDetailByID(pid uint64) (data *model.Post, err error) {
	data = new(model.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id = ?`
	if err = db.Get(data, sqlStr, pid); err != nil {
		zap.L().Error("MySQL执行出错", zap.Error(err), zap.String("sqlStr", sqlStr))
	}
	return
}
