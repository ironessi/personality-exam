// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReconciliationExceptions is the golang structure for table reconciliation_exceptions.
type ReconciliationExceptions struct {
	Id                 uint64      `json:"id"                 orm:"id"                   description:""`                            //
	ReconcileDate      *gtime.Time `json:"reconcileDate"      orm:"reconcile_date"       description:""`                            //
	Channel            string      `json:"channel"            orm:"channel"              description:""`                            //
	ExceptionType      string      `json:"exceptionType"      orm:"exception_type"       description:"LOCAL_MISS/AMOUNT_MISMATCH等"` // LOCAL_MISS/AMOUNT_MISMATCH等
	LocalTransactionNo string      `json:"localTransactionNo" orm:"local_transaction_no" description:""`                            //
	ChannelTradeNo     string      `json:"channelTradeNo"     orm:"channel_trade_no"     description:""`                            //
	OrderNo            string      `json:"orderNo"            orm:"order_no"             description:""`                            //
	LocalAmount        float64     `json:"localAmount"        orm:"local_amount"         description:""`                            //
	ChannelAmount      float64     `json:"channelAmount"      orm:"channel_amount"       description:""`                            //
	LocalStatus        int         `json:"localStatus"        orm:"local_status"         description:""`                            //
	ChannelStatus      string      `json:"channelStatus"      orm:"channel_status"       description:""`                            //
	HandleStatus       int         `json:"handleStatus"       orm:"handle_status"        description:""`                            //
	HandleBy           string      `json:"handleBy"           orm:"handle_by"            description:""`                            //
	HandledAt          *gtime.Time `json:"handledAt"          orm:"handled_at"           description:""`                            //
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:""`                            //
}
