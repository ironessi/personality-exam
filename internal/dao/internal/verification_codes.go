// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// VerificationCodesDao is the data access object for the table verification_codes.
type VerificationCodesDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  VerificationCodesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// VerificationCodesColumns defines and stores column names for the table verification_codes.
type VerificationCodesColumns struct {
	Id        string // 验证码ID
	Mobile    string // 手机号码
	Code      string // 验证码
	Type      string // 业务类型(LOGIN/REGISTER/...)
	Status    string // 状态(0-未使用, 1-已使用, 2-已过期)
	ExpireAt  string // 过期时间
	ClientIp  string // 客户端IP
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// verificationCodesColumns holds the columns for the table verification_codes.
var verificationCodesColumns = VerificationCodesColumns{
	Id:        "id",
	Mobile:    "mobile",
	Code:      "code",
	Type:      "type",
	Status:    "status",
	ExpireAt:  "expire_at",
	ClientIp:  "client_ip",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewVerificationCodesDao creates and returns a new DAO object for table data access.
func NewVerificationCodesDao(handlers ...gdb.ModelHandler) *VerificationCodesDao {
	return &VerificationCodesDao{
		group:    "default",
		table:    "verification_codes",
		columns:  verificationCodesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *VerificationCodesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *VerificationCodesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *VerificationCodesDao) Columns() VerificationCodesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *VerificationCodesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *VerificationCodesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *VerificationCodesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
