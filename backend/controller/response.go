package controller

import (
	"github.com/gin-gonic/gin"
)

type ResCode int16
type ResData struct {
	ResCode `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

const (
	CodeSuccess ResCode = iota + 1000
	CodeNeedToken
	CodeTokenExpired
	CodeInvalidToken
	CodeInvalidParams
	CodeUserExist
	CodeUserPasswordWrong
	CodeContextWrong
	CodeDataNotExist
	CodePageInvalid
	CodeVoteRepeated
	CodeVoteTimeExpired
	CodeServeBusy
)

// Msg decided by ResCode
var codeMap = map[ResCode]string{
	CodeSuccess:           "成功",
	CodeNeedToken:         "需要Token",
	CodeTokenExpired:      "Token已过期",
	CodeInvalidToken:      "无效的Token",
	CodeInvalidParams:     "请求参数错误",
	CodeUserExist:         "用户已存在",
	CodeUserPasswordWrong: "用户密码错误",
	CodeContextWrong:      "上下文中存在问题",
	CodeDataNotExist:      "对应的数据不存在",
	CodePageInvalid:       "页码无效",
	CodeVoteRepeated:      "请勿重复投票",
	CodeVoteTimeExpired:   "投票时间已过",
	CodeServeBusy:         "服务忙",
}

func (c ResCode) Msg() string {
	msg, ok := codeMap[c]
	if !ok {
		return codeMap[CodeServeBusy]
	}
	return msg
}

// Response functions
func ResponseSuccess(c *gin.Context) {
	c.JSON(200, ResData{
		CodeSuccess,
		CodeSuccess.Msg(),
		nil,
	})
}

func ResponseError(c *gin.Context, code int, rcode ResCode) {
	c.JSON(code, ResData{
		rcode,
		rcode.Msg(),
		nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code int, rcode ResCode, msg string) {
	c.JSON(code, ResData{
		rcode,
		msg,
		nil,
	})
}

func ResponseSuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(200, ResData{
		CodeSuccess,
		CodeSuccess.Msg(),
		data,
	})
}
