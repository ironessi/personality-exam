// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReconciliationChannelDataDao is the data access object for the table reconciliation_channel_data.
type ReconciliationChannelDataDao struct {
	table    string                           // table is the underlying table name of the DAO.
	group    string                           // group is the database configuration group name of the current DAO.
	columns  ReconciliationChannelDataColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler               // handlers for customized model modification.
}

// ReconciliationChannelDataColumns defines and stores column names for the table reconciliation_channel_data.
type ReconciliationChannelDataColumns struct {
	Id                 string //
	ReconcileDate      string //
	Channel            string //
	ChannelTradeNo     string //
	ChannelOrderNo     string //
	TradeTime          string //
	Amount             string //
	TradeType          string //
	TradeStatus        string //
	Matched            string //
	LocalTransactionNo string //
	CreatedAt          string //
}

// reconciliationChannelDataColumns holds the columns for the table reconciliation_channel_data.
var reconciliationChannelDataColumns = ReconciliationChannelDataColumns{
	Id:                 "id",
	ReconcileDate:      "reconcile_date",
	Channel:            "channel",
	ChannelTradeNo:     "channel_trade_no",
	ChannelOrderNo:     "channel_order_no",
	TradeTime:          "trade_time",
	Amount:             "amount",
	TradeType:          "trade_type",
	TradeStatus:        "trade_status",
	Matched:            "matched",
	LocalTransactionNo: "local_transaction_no",
	CreatedAt:          "created_at",
}

// NewReconciliationChannelDataDao creates and returns a new DAO object for table data access.
func NewReconciliationChannelDataDao(handlers ...gdb.ModelHandler) *ReconciliationChannelDataDao {
	return &ReconciliationChannelDataDao{
		group:    "default",
		table:    "reconciliation_channel_data",
		columns:  reconciliationChannelDataColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReconciliationChannelDataDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReconciliationChannelDataDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReconciliationChannelDataDao) Columns() ReconciliationChannelDataColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReconciliationChannelDataDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReconciliationChannelDataDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReconciliationChannelDataDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
