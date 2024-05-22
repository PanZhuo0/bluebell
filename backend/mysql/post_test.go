package mysql

import (
	"backend/model"
	"backend/settings"
	"testing"
)

func TestCreatePost(t *testing.T) {
	// 需要一个db对象
	Init(&settings.MySQLConfig{
		Host:         "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "123123",
		DB:           "go",
		MaxOpenConns: 5,
		MaxIdleConns: 10,
	})
	post := model.Post{
		PostID:      10,
		AuthorID:    1,
		CommunityID: 1,
		Title:       "test",
		Content:     "Just a test",
	}
	err := CreatePost(post.AuthorID, &post)
	if err != nil {
		panic(err)
	}
}
