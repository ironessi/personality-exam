// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentQrcodesDao is the data access object for the table payment_qrcodes.
type PaymentQrcodesDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  PaymentQrcodesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// PaymentQrcodesColumns defines and stores column names for the table payment_qrcodes.
type PaymentQrcodesColumns struct {
	Id            string //
	TransactionNo string // 关联支付流水号
	QrType        string // 类型：1-支付宝 2-微信
	QrContent     string // 二维码内容
	QrImageUrl    string // 二维码图片URL
	Status        string // 0-有效 1-已扫描 2-已支付 3-已过期
	ScannedAt     string // 被扫描时间
	ScannedBy     string // 扫描者标识
	ExpireAt      string // 过期时间
	ValidSeconds  string //
	CreatedAt     string //
	UpdatedAt     string //
}

// paymentQrcodesColumns holds the columns for the table payment_qrcodes.
var paymentQrcodesColumns = PaymentQrcodesColumns{
	Id:            "id",
	TransactionNo: "transaction_no",
	QrType:        "qr_type",
	QrContent:     "qr_content",
	QrImageUrl:    "qr_image_url",
	Status:        "status",
	ScannedAt:     "scanned_at",
	ScannedBy:     "scanned_by",
	ExpireAt:      "expire_at",
	ValidSeconds:  "valid_seconds",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewPaymentQrcodesDao creates and returns a new DAO object for table data access.
func NewPaymentQrcodesDao(handlers ...gdb.ModelHandler) *PaymentQrcodesDao {
	return &PaymentQrcodesDao{
		group:    "default",
		table:    "payment_qrcodes",
		columns:  paymentQrcodesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentQrcodesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentQrcodesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentQrcodesDao) Columns() PaymentQrcodesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentQrcodesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentQrcodesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentQrcodesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
