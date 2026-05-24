// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta        `orm:"table:users, do:true"`
	Id            any         // 用户唯一标识，自增主键
	Username      any         // 用户登录名，唯一
	Password      any         // 用户密码，BCrypt加密存储
	Nickname      any         // 用户昵称，用于显示
	AvatarUrl     any         // 用户头像URL地址
	OpenId        any         // 微信开放平台唯一标识
	Mobile        any         // 手机号码
	UserType      any         // 用户类型(0-普通用户,1-管理员)
	CreatedTime   *gtime.Time // 用户创建时间
	LastLoginTime *gtime.Time // 最后登录时间
	Gender        any         // 性别(0-未知,1-男,2-女)
	UpdatedTime   *gtime.Time // 最后更新时间
	DeletedTime   *gtime.Time // 软删除标记
}
