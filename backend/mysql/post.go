package mysql

import (
	"backend/model"
	"backend/utils"

	"go.uber.org/zap"
)

func CreatePost(uid uint64, p *model.Post) (err error) {
	// 1.产生POSTID
	pid, err := utils.GenID()
	if err != nil {
		zap.L().Error("调用utils.GenID()后出错", zap.Error(err))
		return
	}
	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	if _, err = db.Exec(sqlStr, pid, p.Title, p.Content, uid, p.CommunityID); err != nil {
		zap.L().Error("MySQL:执行`insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`时出错", zap.Error(err))
	}
	return
}
