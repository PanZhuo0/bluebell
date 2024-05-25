package controller

import "backend/model"

// 专门用来放接口文档用到的model
// 因为我们的接口文档返回的数据类型是一致的,但是具体的data类型不一样

type _ResponsePostList struct {
	Code    ResCode                `json:"code"`    //业务响应状态码
	Message string                 `json:"message"` //提示信息
	Data    []*model.APIPostDetail `json:"data"`    //数据
}
