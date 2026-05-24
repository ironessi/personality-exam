// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReconciliationChannelData is the golang structure for table reconciliation_channel_data.
type ReconciliationChannelData struct {
	Id                 uint64      `json:"id"                 orm:"id"                   description:""` //
	ReconcileDate      *gtime.Time `json:"reconcileDate"      orm:"reconcile_date"       description:""` //
	Channel            string      `json:"channel"            orm:"channel"              description:""` //
	ChannelTradeNo     string      `json:"channelTradeNo"     orm:"channel_trade_no"     description:""` //
	ChannelOrderNo     string      `json:"channelOrderNo"     orm:"channel_order_no"     description:""` //
	TradeTime          *gtime.Time `json:"tradeTime"          orm:"trade_time"           description:""` //
	Amount             float64     `json:"amount"             orm:"amount"               description:""` //
	TradeType          string      `json:"tradeType"          orm:"trade_type"           description:""` //
	TradeStatus        string      `json:"tradeStatus"        orm:"trade_status"         description:""` //
	Matched            int         `json:"matched"            orm:"matched"              description:""` //
	LocalTransactionNo string      `json:"localTransactionNo" orm:"local_transaction_no" description:""` //
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:""` //
}
