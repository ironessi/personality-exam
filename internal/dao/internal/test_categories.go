// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TestCategoriesDao is the data access object for the table test_categories.
type TestCategoriesDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  TestCategoriesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// TestCategoriesColumns defines and stores column names for the table test_categories.
type TestCategoriesColumns struct {
	Id          string // 分类ID
	Name        string // 分类名称
	Sort        string // 排序优先级
	Status      string // 状态(1-启用, 0-禁用)
	CreatedTime string // 创建时间
	UpdatedTime string // 更新时间
}

// testCategoriesColumns holds the columns for the table test_categories.
var testCategoriesColumns = TestCategoriesColumns{
	Id:          "id",
	Name:        "name",
	Sort:        "sort",
	Status:      "status",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
}

// NewTestCategoriesDao creates and returns a new DAO object for table data access.
func NewTestCategoriesDao(handlers ...gdb.ModelHandler) *TestCategoriesDao {
	return &TestCategoriesDao{
		group:    "default",
		table:    "test_categories",
		columns:  testCategoriesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TestCategoriesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TestCategoriesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TestCategoriesDao) Columns() TestCategoriesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TestCategoriesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TestCategoriesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TestCategoriesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
