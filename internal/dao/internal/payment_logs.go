// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentLogsDao is the data access object for the table payment_logs.
type PaymentLogsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PaymentLogsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PaymentLogsColumns defines and stores column names for the table payment_logs.
type PaymentLogsColumns struct {
	Id            string //
	TransactionNo string //
	OrderNo       string //
	Action        string // 操作类型：CREATE/PREPAY/NOTIFY等
	RequestData   string //
	ResponseData  string //
	Result        string // 0-失败 1-成功
	ErrorMsg      string //
	IpAddress     string //
	CreatedAt     string //
}

// paymentLogsColumns holds the columns for the table payment_logs.
var paymentLogsColumns = PaymentLogsColumns{
	Id:            "id",
	TransactionNo: "transaction_no",
	OrderNo:       "order_no",
	Action:        "action",
	RequestData:   "request_data",
	ResponseData:  "response_data",
	Result:        "result",
	ErrorMsg:      "error_msg",
	IpAddress:     "ip_address",
	CreatedAt:     "created_at",
}

// NewPaymentLogsDao creates and returns a new DAO object for table data access.
func NewPaymentLogsDao(handlers ...gdb.ModelHandler) *PaymentLogsDao {
	return &PaymentLogsDao{
		group:    "default",
		table:    "payment_logs",
		columns:  paymentLogsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentLogsDao) Columns() PaymentLogsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentLogsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
