// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReconciliationExceptionsDao is the data access object for the table reconciliation_exceptions.
type ReconciliationExceptionsDao struct {
	table    string                          // table is the underlying table name of the DAO.
	group    string                          // group is the database configuration group name of the current DAO.
	columns  ReconciliationExceptionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler              // handlers for customized model modification.
}

// ReconciliationExceptionsColumns defines and stores column names for the table reconciliation_exceptions.
type ReconciliationExceptionsColumns struct {
	Id                 string //
	ReconcileDate      string //
	Channel            string //
	ExceptionType      string // LOCAL_MISS/AMOUNT_MISMATCH等
	LocalTransactionNo string //
	ChannelTradeNo     string //
	OrderNo            string //
	LocalAmount        string //
	ChannelAmount      string //
	LocalStatus        string //
	ChannelStatus      string //
	HandleStatus       string //
	HandleBy           string //
	HandledAt          string //
	CreatedAt          string //
}

// reconciliationExceptionsColumns holds the columns for the table reconciliation_exceptions.
var reconciliationExceptionsColumns = ReconciliationExceptionsColumns{
	Id:                 "id",
	ReconcileDate:      "reconcile_date",
	Channel:            "channel",
	ExceptionType:      "exception_type",
	LocalTransactionNo: "local_transaction_no",
	ChannelTradeNo:     "channel_trade_no",
	OrderNo:            "order_no",
	LocalAmount:        "local_amount",
	ChannelAmount:      "channel_amount",
	LocalStatus:        "local_status",
	ChannelStatus:      "channel_status",
	HandleStatus:       "handle_status",
	HandleBy:           "handle_by",
	HandledAt:          "handled_at",
	CreatedAt:          "created_at",
}

// NewReconciliationExceptionsDao creates and returns a new DAO object for table data access.
func NewReconciliationExceptionsDao(handlers ...gdb.ModelHandler) *ReconciliationExceptionsDao {
	return &ReconciliationExceptionsDao{
		group:    "default",
		table:    "reconciliation_exceptions",
		columns:  reconciliationExceptionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReconciliationExceptionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReconciliationExceptionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReconciliationExceptionsDao) Columns() ReconciliationExceptionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReconciliationExceptionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReconciliationExceptionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReconciliationExceptionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
