package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"Auth" summary:"登录"`
	Account  string `v:"required" json:"account" dc:"账号(邮箱/用户名/手机号)"`
	Password string `v:"required" json:"password" dc:"密码"`
}

type LoginRes struct {
	Token string `json:"token"`
}

type RegisterReq struct {
	g.Meta   `path:"/auth/register" method:"post" summary:"注册" tags:"Auth"`
	Mobile   string `v:"required|phone" json:"mobile" dc:"手机号(作为用户名)"`
	Password string `v:"required|length:6,30" json:"password" dc:"密码"`
	Code     string `v:"required" json:"code" dc:"验证码"`
}

type RegisterRes struct{}
