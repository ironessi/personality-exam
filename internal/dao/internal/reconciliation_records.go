// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReconciliationRecordsDao is the data access object for the table reconciliation_records.
type ReconciliationRecordsDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  ReconciliationRecordsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// ReconciliationRecordsColumns defines and stores column names for the table reconciliation_records.
type ReconciliationRecordsColumns struct {
	Id                string //
	ReconcileDate     string // 对账日期
	Channel           string // 渠道
	TotalTransactions string //
	TotalAmount       string //
	MatchedCount      string //
	UnmatchedCount    string //
	Status            string // 0-待对账 1-对账中 2-完成 3-异常
	DetailFileUrl     string //
	ExceptionDetail   string //
	CreatedAt         string //
	CompletedAt       string //
}

// reconciliationRecordsColumns holds the columns for the table reconciliation_records.
var reconciliationRecordsColumns = ReconciliationRecordsColumns{
	Id:                "id",
	ReconcileDate:     "reconcile_date",
	Channel:           "channel",
	TotalTransactions: "total_transactions",
	TotalAmount:       "total_amount",
	MatchedCount:      "matched_count",
	UnmatchedCount:    "unmatched_count",
	Status:            "status",
	DetailFileUrl:     "detail_file_url",
	ExceptionDetail:   "exception_detail",
	CreatedAt:         "created_at",
	CompletedAt:       "completed_at",
}

// NewReconciliationRecordsDao creates and returns a new DAO object for table data access.
func NewReconciliationRecordsDao(handlers ...gdb.ModelHandler) *ReconciliationRecordsDao {
	return &ReconciliationRecordsDao{
		group:    "default",
		table:    "reconciliation_records",
		columns:  reconciliationRecordsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReconciliationRecordsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReconciliationRecordsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReconciliationRecordsDao) Columns() ReconciliationRecordsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReconciliationRecordsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReconciliationRecordsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReconciliationRecordsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
