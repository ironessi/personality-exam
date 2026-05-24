// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FavoritesDao is the data access object for the table favorites.
type FavoritesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  FavoritesColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// FavoritesColumns defines and stores column names for the table favorites.
type FavoritesColumns struct {
	Id          string // 收藏ID，自增主键
	UserId      string // 关联用户ID
	TestTypeId  string // 测试类型ID
	ResultId    string // 测试结果ID（可选）
	CreatedTime string // 收藏创建时间
	UpdatedTime string // 更新时间
	DeletedTime string // 软删除时间
}

// favoritesColumns holds the columns for the table favorites.
var favoritesColumns = FavoritesColumns{
	Id:          "id",
	UserId:      "user_id",
	TestTypeId:  "test_type_id",
	ResultId:    "result_id",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
	DeletedTime: "deleted_time",
}

// NewFavoritesDao creates and returns a new DAO object for table data access.
func NewFavoritesDao(handlers ...gdb.ModelHandler) *FavoritesDao {
	return &FavoritesDao{
		group:    "default",
		table:    "favorites",
		columns:  favoritesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FavoritesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FavoritesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FavoritesDao) Columns() FavoritesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FavoritesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FavoritesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FavoritesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
