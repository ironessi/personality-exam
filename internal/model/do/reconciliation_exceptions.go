// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReconciliationExceptions is the golang structure of table reconciliation_exceptions for DAO operations like Where/Data.
type ReconciliationExceptions struct {
	g.Meta             `orm:"table:reconciliation_exceptions, do:true"`
	Id                 any         //
	ReconcileDate      *gtime.Time //
	Channel            any         //
	ExceptionType      any         // LOCAL_MISS/AMOUNT_MISMATCH等
	LocalTransactionNo any         //
	ChannelTradeNo     any         //
	OrderNo            any         //
	LocalAmount        any         //
	ChannelAmount      any         //
	LocalStatus        any         //
	ChannelStatus      any         //
	HandleStatus       any         //
	HandleBy           any         //
	HandledAt          *gtime.Time //
	CreatedAt          *gtime.Time //
}
