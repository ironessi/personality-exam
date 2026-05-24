// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RefundRecords is the golang structure of table refund_records for DAO operations like Where/Data.
type RefundRecords struct {
	g.Meta          `orm:"table:refund_records, do:true"`
	Id              any         //
	RefundNo        any         // 内部退款单号
	TransactionNo   any         // 原支付流水号
	OrderNo         any         // 原订单号
	RefundAmount    any         // 退款金额
	Reason          any         // 退款原因
	Status          any         // 0-申请中 1-处理中 2-成功 3-失败
	ChannelRefundId any         // 渠道退款单号
	ChannelResponse any         //
	CreatedAt       *gtime.Time //
	CompletedAt     *gtime.Time //
}
