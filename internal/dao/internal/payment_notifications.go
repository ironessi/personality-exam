// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentNotificationsDao is the data access object for the table payment_notifications.
type PaymentNotificationsDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  PaymentNotificationsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// PaymentNotificationsColumns defines and stores column names for the table payment_notifications.
type PaymentNotificationsColumns struct {
	Id            string //
	Channel       string // alipay/wechat
	NotifyType    string // 类型：pay/refund/dispute
	RawHeaders    string // HTTP头
	RawBody       string // 原始请求体
	ParsedData    string // 解析后的数据
	Status        string // 0-待处理 1-处理中 2-成功 3-失败
	RetryCount    string //
	ErrorMsg      string //
	TransactionNo string //
	OrderNo       string //
	ReceivedAt    string //
	ProcessedAt   string //
}

// paymentNotificationsColumns holds the columns for the table payment_notifications.
var paymentNotificationsColumns = PaymentNotificationsColumns{
	Id:            "id",
	Channel:       "channel",
	NotifyType:    "notify_type",
	RawHeaders:    "raw_headers",
	RawBody:       "raw_body",
	ParsedData:    "parsed_data",
	Status:        "status",
	RetryCount:    "retry_count",
	ErrorMsg:      "error_msg",
	TransactionNo: "transaction_no",
	OrderNo:       "order_no",
	ReceivedAt:    "received_at",
	ProcessedAt:   "processed_at",
}

// NewPaymentNotificationsDao creates and returns a new DAO object for table data access.
func NewPaymentNotificationsDao(handlers ...gdb.ModelHandler) *PaymentNotificationsDao {
	return &PaymentNotificationsDao{
		group:    "default",
		table:    "payment_notifications",
		columns:  paymentNotificationsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentNotificationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentNotificationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentNotificationsDao) Columns() PaymentNotificationsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentNotificationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentNotificationsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentNotificationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
