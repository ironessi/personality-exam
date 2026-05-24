// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentNotifications is the golang structure of table payment_notifications for DAO operations like Where/Data.
type PaymentNotifications struct {
	g.Meta        `orm:"table:payment_notifications, do:true"`
	Id            any         //
	Channel       any         // alipay/wechat
	NotifyType    any         // 类型：pay/refund/dispute
	RawHeaders    any         // HTTP头
	RawBody       any         // 原始请求体
	ParsedData    any         // 解析后的数据
	Status        any         // 0-待处理 1-处理中 2-成功 3-失败
	RetryCount    any         //
	ErrorMsg      any         //
	TransactionNo any         //
	OrderNo       any         //
	ReceivedAt    *gtime.Time //
	ProcessedAt   *gtime.Time //
}
