// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentLogs is the golang structure for table payment_logs.
type PaymentLogs struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`                           //
	TransactionNo string      `json:"transactionNo" orm:"transaction_no" description:""`                           //
	OrderNo       string      `json:"orderNo"       orm:"order_no"       description:""`                           //
	Action        string      `json:"action"        orm:"action"         description:"操作类型：CREATE/PREPAY/NOTIFY等"` // 操作类型：CREATE/PREPAY/NOTIFY等
	RequestData   string      `json:"requestData"   orm:"request_data"   description:""`                           //
	ResponseData  string      `json:"responseData"  orm:"response_data"  description:""`                           //
	Result        int         `json:"result"        orm:"result"         description:"0-失败 1-成功"`                  // 0-失败 1-成功
	ErrorMsg      string      `json:"errorMsg"      orm:"error_msg"      description:""`                           //
	IpAddress     string      `json:"ipAddress"     orm:"ip_address"     description:""`                           //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`                           //
}
