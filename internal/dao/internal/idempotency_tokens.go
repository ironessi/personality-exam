// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// IdempotencyTokensDao is the data access object for the table idempotency_tokens.
type IdempotencyTokensDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  IdempotencyTokensColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// IdempotencyTokensColumns defines and stores column names for the table idempotency_tokens.
type IdempotencyTokensColumns struct {
	Token        string //
	BusinessType string // 业务类型
	BusinessId   string // 业务单号
	Status       string // 0-处理中 1-成功 2-失败
	ExpireAt     string //
	CreatedAt    string //
}

// idempotencyTokensColumns holds the columns for the table idempotency_tokens.
var idempotencyTokensColumns = IdempotencyTokensColumns{
	Token:        "token",
	BusinessType: "business_type",
	BusinessId:   "business_id",
	Status:       "status",
	ExpireAt:     "expire_at",
	CreatedAt:    "created_at",
}

// NewIdempotencyTokensDao creates and returns a new DAO object for table data access.
func NewIdempotencyTokensDao(handlers ...gdb.ModelHandler) *IdempotencyTokensDao {
	return &IdempotencyTokensDao{
		group:    "default",
		table:    "idempotency_tokens",
		columns:  idempotencyTokensColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *IdempotencyTokensDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *IdempotencyTokensDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *IdempotencyTokensDao) Columns() IdempotencyTokensColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *IdempotencyTokensDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *IdempotencyTokensDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *IdempotencyTokensDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
