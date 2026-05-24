// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserTestHistoryDao is the data access object for the table user_test_history.
type UserTestHistoryDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  UserTestHistoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// UserTestHistoryColumns defines and stores column names for the table user_test_history.
type UserTestHistoryColumns struct {
	Id          string // 历史记录ID，自增主键
	UserId      string // 关联用户ID
	TestTypeId  string // 测试类型ID
	ResultId    string // 关联的测试结果ID
	ResultType  string // 结果类型代码
	ResultTitle string // 结果标题
	WingType    string // 翼型(九型人格专用)
	CreateTime  string // 创建时间
}

// userTestHistoryColumns holds the columns for the table user_test_history.
var userTestHistoryColumns = UserTestHistoryColumns{
	Id:          "id",
	UserId:      "user_id",
	TestTypeId:  "test_type_id",
	ResultId:    "result_id",
	ResultType:  "result_type",
	ResultTitle: "result_title",
	WingType:    "wing_type",
	CreateTime:  "create_time",
}

// NewUserTestHistoryDao creates and returns a new DAO object for table data access.
func NewUserTestHistoryDao(handlers ...gdb.ModelHandler) *UserTestHistoryDao {
	return &UserTestHistoryDao{
		group:    "default",
		table:    "user_test_history",
		columns:  userTestHistoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserTestHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserTestHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserTestHistoryDao) Columns() UserTestHistoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserTestHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserTestHistoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserTestHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
