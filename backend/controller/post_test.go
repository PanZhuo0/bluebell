package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 帖子的创建Controller层测试
func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode) //测试模式的GIN
	r := gin.Default()
	url := `/api/v1/post`
	r.POST(url, CreatePostHandler)

	body := `{
		"community_id":10,
		"title":"Test",
		"content":"just a test"
	}`
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json") //需要在请求头中标记文件类型为json
	w := httptest.NewRecorder()                        //记录请求的结果
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "上下文中存在问题") //结果是否为需要登录Token？
}

// 查询帖子测试
// 按ID获取帖子详情测试
