// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TestOrdersDao is the data access object for the table test_orders.
type TestOrdersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TestOrdersColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TestOrdersColumns defines and stores column names for the table test_orders.
type TestOrdersColumns struct {
	Id                string // 订单ID
	OrderNo           string // 订单编号（业务唯一）
	UserId            string // 用户ID
	TestId            string // 测试类型ID
	TestName          string // 测试名称（冗余存储）
	OriginalPrice     string // 原价
	DiscountAmount    string // 优惠金额
	PayableAmount     string // 应付金额（实际支付）
	Status            string // 状态：0-待支付 1-已支付 2-处理中 3-已完成 4-已退款 5-已关闭
	PaidAt            string // 支付完成时间
	PaidChannel       string // 支付渠道：alipay/wechat
	PaidTransactionId string // 第三方支付流水号
	ResultReportId    string // 生成的报告ID
	ExpireAt          string // 订单过期时间
	ClientIp          string // 用户下单IP
	UserAgent         string // 用户代理
	IsDeleted         string //
	CreatedAt         string //
	UpdatedAt         string //
}

// testOrdersColumns holds the columns for the table test_orders.
var testOrdersColumns = TestOrdersColumns{
	Id:                "id",
	OrderNo:           "order_no",
	UserId:            "user_id",
	TestId:            "test_id",
	TestName:          "test_name",
	OriginalPrice:     "original_price",
	DiscountAmount:    "discount_amount",
	PayableAmount:     "payable_amount",
	Status:            "status",
	PaidAt:            "paid_at",
	PaidChannel:       "paid_channel",
	PaidTransactionId: "paid_transaction_id",
	ResultReportId:    "result_report_id",
	ExpireAt:          "expire_at",
	ClientIp:          "client_ip",
	UserAgent:         "user_agent",
	IsDeleted:         "is_deleted",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewTestOrdersDao creates and returns a new DAO object for table data access.
func NewTestOrdersDao(handlers ...gdb.ModelHandler) *TestOrdersDao {
	return &TestOrdersDao{
		group:    "default",
		table:    "test_orders",
		columns:  testOrdersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TestOrdersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TestOrdersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TestOrdersDao) Columns() TestOrdersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TestOrdersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TestOrdersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TestOrdersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
