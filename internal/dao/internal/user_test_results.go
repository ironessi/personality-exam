// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserTestResultsDao is the data access object for the table user_test_results.
type UserTestResultsDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  UserTestResultsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// UserTestResultsColumns defines and stores column names for the table user_test_results.
type UserTestResultsColumns struct {
	Id          string // 结果ID，自增主键
	UserId      string // 用户ID
	TestTypeId  string // 测试类型ID
	ResultData  string // 结果数据(包含总分、分析结果、建议等)
	TestDate    string // 测试时间
	CreatedTime string // 创建时间
}

// userTestResultsColumns holds the columns for the table user_test_results.
var userTestResultsColumns = UserTestResultsColumns{
	Id:          "id",
	UserId:      "user_id",
	TestTypeId:  "test_type_id",
	ResultData:  "result_data",
	TestDate:    "test_date",
	CreatedTime: "created_time",
}

// NewUserTestResultsDao creates and returns a new DAO object for table data access.
func NewUserTestResultsDao(handlers ...gdb.ModelHandler) *UserTestResultsDao {
	return &UserTestResultsDao{
		group:    "default",
		table:    "user_test_results",
		columns:  userTestResultsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserTestResultsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserTestResultsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserTestResultsDao) Columns() UserTestResultsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserTestResultsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserTestResultsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserTestResultsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
