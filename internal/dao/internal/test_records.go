// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TestRecordsDao is the data access object for the table test_records.
type TestRecordsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TestRecordsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TestRecordsColumns defines and stores column names for the table test_records.
type TestRecordsColumns struct {
	Id          string // 记录ID，自增主键
	UserId      string // 关联用户ID
	TestTypeId  string // 测试类型ID
	Status      string // 测试状态(进行中/已完成)
	StartTime   string // 测试开始时间
	EndTime     string // 测试结束时间
	CreatedTime string // 记录创建时间
	UpdatedTime string // 记录更新时间
	Deleted     string // 软删除标记(0-未删除,1-已删除)
}

// testRecordsColumns holds the columns for the table test_records.
var testRecordsColumns = TestRecordsColumns{
	Id:          "id",
	UserId:      "user_id",
	TestTypeId:  "test_type_id",
	Status:      "status",
	StartTime:   "start_time",
	EndTime:     "end_time",
	CreatedTime: "created_time",
	UpdatedTime: "updated_time",
	Deleted:     "deleted",
}

// NewTestRecordsDao creates and returns a new DAO object for table data access.
func NewTestRecordsDao(handlers ...gdb.ModelHandler) *TestRecordsDao {
	return &TestRecordsDao{
		group:    "default",
		table:    "test_records",
		columns:  testRecordsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TestRecordsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TestRecordsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TestRecordsDao) Columns() TestRecordsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TestRecordsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TestRecordsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TestRecordsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
