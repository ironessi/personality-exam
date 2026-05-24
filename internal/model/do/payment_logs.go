// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentLogs is the golang structure of table payment_logs for DAO operations like Where/Data.
type PaymentLogs struct {
	g.Meta        `orm:"table:payment_logs, do:true"`
	Id            any         //
	TransactionNo any         //
	OrderNo       any         //
	Action        any         // 操作类型：CREATE/PREPAY/NOTIFY等
	RequestData   any         //
	ResponseData  any         //
	Result        any         // 0-失败 1-成功
	ErrorMsg      any         //
	IpAddress     any         //
	CreatedAt     *gtime.Time //
}
