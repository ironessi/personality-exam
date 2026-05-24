// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id            int64       `json:"id"            orm:"id"              description:"用户唯一标识，自增主键"`        // 用户唯一标识，自增主键
	Username      string      `json:"username"      orm:"username"        description:"用户登录名，唯一"`           // 用户登录名，唯一
	Password      string      `json:"password"      orm:"password"        description:"用户密码，BCrypt加密存储"`    // 用户密码，BCrypt加密存储
	Nickname      string      `json:"nickname"      orm:"nickname"        description:"用户昵称，用于显示"`          // 用户昵称，用于显示
	AvatarUrl     string      `json:"avatarUrl"     orm:"avatar_url"      description:"用户头像URL地址"`          // 用户头像URL地址
	OpenId        string      `json:"openId"        orm:"open_id"         description:"微信开放平台唯一标识"`         // 微信开放平台唯一标识
	Mobile        string      `json:"mobile"        orm:"mobile"          description:"手机号码"`               // 手机号码
	UserType      int         `json:"userType"      orm:"user_type"       description:"用户类型(0-普通用户,1-管理员)"` // 用户类型(0-普通用户,1-管理员)
	CreatedTime   *gtime.Time `json:"createdTime"   orm:"created_time"    description:"用户创建时间"`             // 用户创建时间
	LastLoginTime *gtime.Time `json:"lastLoginTime" orm:"last_login_time" description:"最后登录时间"`             // 最后登录时间
	Gender        int         `json:"gender"        orm:"gender"          description:"性别(0-未知,1-男,2-女)"`   // 性别(0-未知,1-男,2-女)
	UpdatedTime   *gtime.Time `json:"updatedTime"   orm:"updated_time"    description:"最后更新时间"`             // 最后更新时间
	DeletedTime   *gtime.Time `json:"deletedTime"   orm:"deleted_time"    description:"软删除标记"`              // 软删除标记
}
