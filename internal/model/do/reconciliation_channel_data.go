// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReconciliationChannelData is the golang structure of table reconciliation_channel_data for DAO operations like Where/Data.
type ReconciliationChannelData struct {
	g.Meta             `orm:"table:reconciliation_channel_data, do:true"`
	Id                 any         //
	ReconcileDate      *gtime.Time //
	Channel            any         //
	ChannelTradeNo     any         //
	ChannelOrderNo     any         //
	TradeTime          *gtime.Time //
	Amount             any         //
	TradeType          any         //
	TradeStatus        any         //
	Matched            any         //
	LocalTransactionNo any         //
	CreatedAt          *gtime.Time //
}
