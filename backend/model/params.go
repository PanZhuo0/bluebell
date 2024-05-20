package model

// 参数
type ParamSignUp struct {
	UserName   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type VoteData struct {
	PostID    uint64  `json:"post_id" binding:"required"`
	Direction float64 `json:"direction" binding:"required"`
}
