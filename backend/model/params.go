package model

import (
	"encoding/json"
	"errors"
)

// 参数
type ParamSignUp struct {
	UserName   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type VoteData struct {
	PostID    uint64  `json:"post_id"`
	Direction float64 `json:"direction"`
}

func (v *VoteData) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		PostID    uint64  `json:"post_id"`
		Direction float64 `json:"direction"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if (required.PostID) == 0 {
		err = errors.New("缺少必填字段post_id")
	} else if required.Direction == 0 {
		err = errors.New("缺少必填字段direction")
	} else if required.Direction != 0 && required.Direction != -1 && required.Direction != 1 {
		err = errors.New("direction 只能为1/-1")
	} else {
		v.PostID = required.PostID
		v.Direction = required.Direction
	}
	return
}

const OrderTime = "time"
const OrderScore = "score"

type ParamPostList struct {
	CommunityID uint64 `form:"community_id"`
	Page        int    `form:"page"`
	Size        int    `form:"size"`
	Order       string `form:"order"`
}
