package model

import "github.com/gogf/gf/v2/os/gtime"

// SmsSendReq 发送验证码请求
type SmsSendReq struct {
	Mobile       string `json:"mobile" v:"required|phone#手机号不能为空|手机号格式不正确"`
	Type         string `json:"type" v:"required#验证码类型不能为空"` // LOGIN, REGISTER, FORGOT
	TemplateCode string `json:"templateCode"`                // 可选，默认使用配置
}

// SmsSendRes 发送验证码响应
type SmsSendRes struct {
	RequestId string `json:"requestId"`
	ExpireIn  int    `json:"expireIn"` // 过期时间(秒)
}

// SmsVerifyReq 验证验证码请求
type SmsVerifyReq struct {
	Mobile string `json:"mobile" v:"required|phone#手机号不能为空|手机号格式不正确"`
	Code   string `json:"code" v:"required#验证码不能为空"`
	Type   string `json:"type" v:"required#验证类型不能为空"`
}

// SmsVerifyRes 验证验证码响应
type SmsVerifyRes struct {
	Token string      `json:"token,omitempty"`
	User  *UserDetail `json:"user,omitempty"`
}

// UserDetail 用户详情 (复用或简化 User Entity)
type UserDetail struct {
	Id       int64  `json:"id"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
}

// SmsStatusReq 查询状态请求
type SmsStatusReq struct {
	Mobile string `json:"mobile" v:"required|phone#手机号不能为空|手机号格式不正确"`
}

// SmsStatusRes 查询状态响应
type SmsStatusRes struct {
	HasSent           bool        `json:"hasSent"`
	LastSendTime      *gtime.Time `json:"lastSendTime"`
	RemainingCount    int         `json:"remainingCount"`
	NextAvailableTime *gtime.Time `json:"nextAvailableTime"`
}
