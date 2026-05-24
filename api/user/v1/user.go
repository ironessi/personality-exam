package v1

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type GetUserInfoReq struct {
	g.Meta `path:"user/info" method:"get" summary:"获取用户信息" tags:"用户"`
	Id     int64 `json:"id" v:"required" json:"id" dc:"用户ID"`
}

type GetUserInfoRes struct {
	Id            int64     `json:"id"`       // 用户唯一标识，自增主键
	Username      string    `json:"username"` // 用户登录名，唯一
	Nickname      string    `json:"nickname"` // 用户昵称，用于显示
	Gender        int       `json:"gender"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}

type UpdateUserInfoReq struct {
	g.Meta   `path:"user/update" method:"post" summary:"更新用户信息" tags:"用户"`
	Id       int64  `json:"id" v:"required" json:"id" dc:"用户ID"`
	Nickname string `json:"nickname" dc:"昵称"`
	Gender   int    `json:"gender" dc:"性别"`
}

type UpdateUserInfoRes struct{}
