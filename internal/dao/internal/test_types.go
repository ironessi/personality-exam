// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TestTypesDao is the data access object for the table test_types.
type TestTypesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TestTypesColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TestTypesColumns defines and stores column names for the table test_types.
type TestTypesColumns struct {
	Id          string // 类型ID，自增主键
	Name        string // 测试名称
	Description string // 测试描述
	TestTime    string // 预计测试时间(分钟)
	Status      string // 状态(0-禁用,1-启用)
	CreatedTime string // 创建时间
	UpdatedTime string // 更新时间
	Img         string // 测试图标URL
	CategoryId  string // 所属分类ID
	Price       string // 查看报告所需支付金额（元）
}

// testTypesColumns holds the columns for the table test_types.
var testTypesColumns = TestTypesColumns{
	Id:          "id",
	Name:        "name",
	Description: "description",
	TestTime:    "test_time",
	Status:      "status",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
	Img:         "img",
	CategoryId:  "category_id",
	Price:       "price",
}

// NewTestTypesDao creates and returns a new DAO object for table data access.
func NewTestTypesDao(handlers ...gdb.ModelHandler) *TestTypesDao {
	return &TestTypesDao{
		group:    "default",
		table:    "test_types",
		columns:  testTypesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TestTypesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TestTypesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TestTypesDao) Columns() TestTypesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TestTypesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TestTypesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TestTypesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
