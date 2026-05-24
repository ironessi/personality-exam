// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RefundRecords is the golang structure for table refund_records.
type RefundRecords struct {
	Id              uint64      `json:"id"              orm:"id"                description:""`                      //
	RefundNo        string      `json:"refundNo"        orm:"refund_no"         description:"内部退款单号"`                // 内部退款单号
	TransactionNo   string      `json:"transactionNo"   orm:"transaction_no"    description:"原支付流水号"`                // 原支付流水号
	OrderNo         string      `json:"orderNo"         orm:"order_no"          description:"原订单号"`                  // 原订单号
	RefundAmount    float64     `json:"refundAmount"    orm:"refund_amount"     description:"退款金额"`                  // 退款金额
	Reason          string      `json:"reason"          orm:"reason"            description:"退款原因"`                  // 退款原因
	Status          int         `json:"status"          orm:"status"            description:"0-申请中 1-处理中 2-成功 3-失败"` // 0-申请中 1-处理中 2-成功 3-失败
	ChannelRefundId string      `json:"channelRefundId" orm:"channel_refund_id" description:"渠道退款单号"`                // 渠道退款单号
	ChannelResponse string      `json:"channelResponse" orm:"channel_response"  description:""`                      //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        description:""`                      //
	CompletedAt     *gtime.Time `json:"completedAt"     orm:"completed_at"      description:""`                      //
}
