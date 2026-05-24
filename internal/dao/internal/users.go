// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersDao is the data access object for the table users.
type UsersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UsersColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UsersColumns defines and stores column names for the table users.
type UsersColumns struct {
	Id            string // 用户唯一标识，自增主键
	Username      string // 用户登录名，唯一
	Password      string // 用户密码，BCrypt加密存储
	Nickname      string // 用户昵称，用于显示
	AvatarUrl     string // 用户头像URL地址
	OpenId        string // 微信开放平台唯一标识
	Mobile        string // 手机号码
	UserType      string // 用户类型(0-普通用户,1-管理员)
	CreatedTime   string // 用户创建时间
	LastLoginTime string // 最后登录时间
	Gender        string // 性别(0-未知,1-男,2-女)
	UpdatedTime   string // 最后更新时间
	DeletedTime   string // 软删除标记
}

// usersColumns holds the columns for the table users.
var usersColumns = UsersColumns{
	Id:            "id",
	Username:      "username",
	Password:      "password",
	Nickname:      "nickname",
	AvatarUrl:     "avatar_url",
	OpenId:        "open_id",
	Mobile:        "mobile",
	UserType:      "user_type",
	CreatedTime:   "created_time",
	LastLoginTime: "last_login_time",
	Gender:        "gender",
	UpdatedTime:   "updated_time",
	DeletedTime:   "deleted_time",
}

// NewUsersDao creates and returns a new DAO object for table data access.
func NewUsersDao(handlers ...gdb.ModelHandler) *UsersDao {
	return &UsersDao{
		group:    "default",
		table:    "users",
		columns:  usersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UsersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UsersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UsersDao) Columns() UsersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UsersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UsersDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *UsersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
