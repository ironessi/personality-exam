// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RefundRecordsDao is the data access object for the table refund_records.
type RefundRecordsDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  RefundRecordsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// RefundRecordsColumns defines and stores column names for the table refund_records.
type RefundRecordsColumns struct {
	Id              string //
	RefundNo        string // 内部退款单号
	TransactionNo   string // 原支付流水号
	OrderNo         string // 原订单号
	RefundAmount    string // 退款金额
	Reason          string // 退款原因
	Status          string // 0-申请中 1-处理中 2-成功 3-失败
	ChannelRefundId string // 渠道退款单号
	ChannelResponse string //
	CreatedAt       string //
	CompletedAt     string //
}

// refundRecordsColumns holds the columns for the table refund_records.
var refundRecordsColumns = RefundRecordsColumns{
	Id:              "id",
	RefundNo:        "refund_no",
	TransactionNo:   "transaction_no",
	OrderNo:         "order_no",
	RefundAmount:    "refund_amount",
	Reason:          "reason",
	Status:          "status",
	ChannelRefundId: "channel_refund_id",
	ChannelResponse: "channel_response",
	CreatedAt:       "created_at",
	CompletedAt:     "completed_at",
}

// NewRefundRecordsDao creates and returns a new DAO object for table data access.
func NewRefundRecordsDao(handlers ...gdb.ModelHandler) *RefundRecordsDao {
	return &RefundRecordsDao{
		group:    "default",
		table:    "refund_records",
		columns:  refundRecordsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RefundRecordsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RefundRecordsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RefundRecordsDao) Columns() RefundRecordsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RefundRecordsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RefundRecordsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RefundRecordsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
