// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTransactions is the golang structure of table payment_transactions for DAO operations like Where/Data.
type PaymentTransactions struct {
	g.Meta              `orm:"table:payment_transactions, do:true"`
	Id                  any         //
	TransactionNo       any         // 内部支付流水号（唯一）
	OrderNo             any         // 业务订单号
	UserId              any         // 用户ID
	Amount              any         // 支付金额
	Currency            any         // 币种
	Channel             any         // 支付渠道：alipay/wechat
	ChannelAppId        any         // 渠道应用ID
	Status              any         // 状态：0-待支付 1-支付中 2-支付成功 3-支付失败 4-已关闭 5-已退款
	ChannelExtra        any         // 渠道扩展参数
	AlipayTradeNo       any         // 支付宝交易号
	WechatPrepayId      any         // 微信预支付ID
	WechatTransactionId any         // 微信订单号
	PrepayAt            any         // 预下单时间
	PaidAt              *gtime.Time // 支付成功时间
	ClosedAt            *gtime.Time // 关闭时间
	ExpireAt            any         // 过期时间
	NotifyRawData       any         // 原始回调数据
	NotifiedAt          *gtime.Time // 收到回调时间
	IdempotencyKey      any         // 幂等键
	ClientIp            any         //
	CreatedAt           *gtime.Time //
	UpdatedAt           *gtime.Time //
}
