package mysql

import (
	"backend/model"
	"strings"

	"github.com/jmoiron/sqlx"
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

// 依据从Redis中取出的IDs从MySQL中取出对应的数据
func GetPostListByIDs(ids []string) (data []*model.Post, err error) {
	sqlStr := `
		select 
			post_id,title,content,author_id,community_id,create_time
		from
			post
		where 
			post_id 
		in 
			(?)
		order by
			FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ",")) //第一个ids 放在 in (?) 第二个ids.join(,) 放在 FIND_IN_SET(post_id,?)上
	if err != nil {
		zap.L().Error("调用sqlx.In()后出错", zap.String("sqlStr", sqlStr), zap.Any("ids[]", ids))
		return
	}
	query = db.Rebind(query) //重新绑定？
	if err = db.Select(&data, query, args...); err != nil {
		zap.L().Error("MySQL:执行批量查询失败", zap.Error(err), zap.String("query", query))
		return
	}
	return
}
