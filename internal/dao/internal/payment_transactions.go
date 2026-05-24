// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentTransactionsDao is the data access object for the table payment_transactions.
type PaymentTransactionsDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  PaymentTransactionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// PaymentTransactionsColumns defines and stores column names for the table payment_transactions.
type PaymentTransactionsColumns struct {
	Id                  string //
	TransactionNo       string // 内部支付流水号（唯一）
	OrderNo             string // 业务订单号
	UserId              string // 用户ID
	Amount              string // 支付金额
	Currency            string // 币种
	Channel             string // 支付渠道：alipay/wechat
	ChannelAppId        string // 渠道应用ID
	Status              string // 状态：0-待支付 1-支付中 2-支付成功 3-支付失败 4-已关闭 5-已退款
	ChannelExtra        string // 渠道扩展参数
	AlipayTradeNo       string // 支付宝交易号
	WechatPrepayId      string // 微信预支付ID
	WechatTransactionId string // 微信订单号
	PrepayAt            string // 预下单时间
	PaidAt              string // 支付成功时间
	ClosedAt            string // 关闭时间
	ExpireAt            string // 过期时间
	NotifyRawData       string // 原始回调数据
	NotifiedAt          string // 收到回调时间
	IdempotencyKey      string // 幂等键
	ClientIp            string //
	CreatedAt           string //
	UpdatedAt           string //
}

// paymentTransactionsColumns holds the columns for the table payment_transactions.
var paymentTransactionsColumns = PaymentTransactionsColumns{
	Id:                  "id",
	TransactionNo:       "transaction_no",
	OrderNo:             "order_no",
	UserId:              "user_id",
	Amount:              "amount",
	Currency:            "currency",
	Channel:             "channel",
	ChannelAppId:        "channel_app_id",
	Status:              "status",
	ChannelExtra:        "channel_extra",
	AlipayTradeNo:       "alipay_trade_no",
	WechatPrepayId:      "wechat_prepay_id",
	WechatTransactionId: "wechat_transaction_id",
	PrepayAt:            "prepay_at",
	PaidAt:              "paid_at",
	ClosedAt:            "closed_at",
	ExpireAt:            "expire_at",
	NotifyRawData:       "notify_raw_data",
	NotifiedAt:          "notified_at",
	IdempotencyKey:      "idempotency_key",
	ClientIp:            "client_ip",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewPaymentTransactionsDao creates and returns a new DAO object for table data access.
func NewPaymentTransactionsDao(handlers ...gdb.ModelHandler) *PaymentTransactionsDao {
	return &PaymentTransactionsDao{
		group:    "default",
		table:    "payment_transactions",
		columns:  paymentTransactionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PaymentTransactionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PaymentTransactionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PaymentTransactionsDao) Columns() PaymentTransactionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PaymentTransactionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PaymentTransactionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PaymentTransactionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
