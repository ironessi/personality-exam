// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentTransactions is the golang structure for table payment_transactions.
type PaymentTransactions struct {
	Id                  uint64      `json:"id"                  orm:"id"                    description:""`                                         //
	TransactionNo       string      `json:"transactionNo"       orm:"transaction_no"        description:"内部支付流水号（唯一）"`                              // 内部支付流水号（唯一）
	OrderNo             string      `json:"orderNo"             orm:"order_no"              description:"业务订单号"`                                    // 业务订单号
	UserId              uint64      `json:"userId"              orm:"user_id"               description:"用户ID"`                                     // 用户ID
	Amount              float64     `json:"amount"              orm:"amount"                description:"支付金额"`                                     // 支付金额
	Currency            string      `json:"currency"            orm:"currency"              description:"币种"`                                       // 币种
	Channel             string      `json:"channel"             orm:"channel"               description:"支付渠道：alipay/wechat"`                       // 支付渠道：alipay/wechat
	ChannelAppId        string      `json:"channelAppId"        orm:"channel_app_id"        description:"渠道应用ID"`                                   // 渠道应用ID
	Status              int         `json:"status"              orm:"status"                description:"状态：0-待支付 1-支付中 2-支付成功 3-支付失败 4-已关闭 5-已退款"` // 状态：0-待支付 1-支付中 2-支付成功 3-支付失败 4-已关闭 5-已退款
	ChannelExtra        string      `json:"channelExtra"        orm:"channel_extra"         description:"渠道扩展参数"`                                   // 渠道扩展参数
	AlipayTradeNo       string      `json:"alipayTradeNo"       orm:"alipay_trade_no"       description:"支付宝交易号"`                                   // 支付宝交易号
	WechatPrepayId      string      `json:"wechatPrepayId"      orm:"wechat_prepay_id"      description:"微信预支付ID"`                                  // 微信预支付ID
	WechatTransactionId string      `json:"wechatTransactionId" orm:"wechat_transaction_id" description:"微信订单号"`                                    // 微信订单号
	PrepayAt            int64       `json:"prepayAt"            orm:"prepay_at"             description:"预下单时间"`                                    // 预下单时间
	PaidAt              *gtime.Time `json:"paidAt"              orm:"paid_at"               description:"支付成功时间"`                                   // 支付成功时间
	ClosedAt            *gtime.Time `json:"closedAt"            orm:"closed_at"             description:"关闭时间"`                                     // 关闭时间
	ExpireAt            int64       `json:"expireAt"            orm:"expire_at"             description:"过期时间"`                                     // 过期时间
	NotifyRawData       string      `json:"notifyRawData"       orm:"notify_raw_data"       description:"原始回调数据"`                                   // 原始回调数据
	NotifiedAt          *gtime.Time `json:"notifiedAt"          orm:"notified_at"           description:"收到回调时间"`                                   // 收到回调时间
	IdempotencyKey      string      `json:"idempotencyKey"      orm:"idempotency_key"       description:"幂等键"`                                      // 幂等键
	ClientIp            string      `json:"clientIp"            orm:"client_ip"             description:""`                                         //
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            description:""`                                         //
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            description:""`                                         //
}
