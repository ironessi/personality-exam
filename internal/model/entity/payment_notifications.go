// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentNotifications is the golang structure for table payment_notifications.
type PaymentNotifications struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`                      //
	Channel       string      `json:"channel"       orm:"channel"        description:"alipay/wechat"`         // alipay/wechat
	NotifyType    string      `json:"notifyType"    orm:"notify_type"    description:"类型：pay/refund/dispute"` // 类型：pay/refund/dispute
	RawHeaders    string      `json:"rawHeaders"    orm:"raw_headers"    description:"HTTP头"`                 // HTTP头
	RawBody       string      `json:"rawBody"       orm:"raw_body"       description:"原始请求体"`                 // 原始请求体
	ParsedData    string      `json:"parsedData"    orm:"parsed_data"    description:"解析后的数据"`                // 解析后的数据
	Status        int         `json:"status"        orm:"status"         description:"0-待处理 1-处理中 2-成功 3-失败"` // 0-待处理 1-处理中 2-成功 3-失败
	RetryCount    uint        `json:"retryCount"    orm:"retry_count"    description:""`                      //
	ErrorMsg      string      `json:"errorMsg"      orm:"error_msg"      description:""`                      //
	TransactionNo string      `json:"transactionNo" orm:"transaction_no" description:""`                      //
	OrderNo       string      `json:"orderNo"       orm:"order_no"       description:""`                      //
	ReceivedAt    *gtime.Time `json:"receivedAt"    orm:"received_at"    description:""`                      //
	ProcessedAt   *gtime.Time `json:"processedAt"   orm:"processed_at"   description:""`                      //
}
