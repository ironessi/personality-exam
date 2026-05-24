// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ReconciliationRecords is the golang structure for table reconciliation_records.
type ReconciliationRecords struct {
	Id                uint64      `json:"id"                orm:"id"                 description:""`                      //
	ReconcileDate     *gtime.Time `json:"reconcileDate"     orm:"reconcile_date"     description:"对账日期"`                  // 对账日期
	Channel           string      `json:"channel"           orm:"channel"            description:"渠道"`                    // 渠道
	TotalTransactions uint        `json:"totalTransactions" orm:"total_transactions" description:""`                      //
	TotalAmount       float64     `json:"totalAmount"       orm:"total_amount"       description:""`                      //
	MatchedCount      uint        `json:"matchedCount"      orm:"matched_count"      description:""`                      //
	UnmatchedCount    uint        `json:"unmatchedCount"    orm:"unmatched_count"    description:""`                      //
	Status            int         `json:"status"            orm:"status"             description:"0-待对账 1-对账中 2-完成 3-异常"` // 0-待对账 1-对账中 2-完成 3-异常
	DetailFileUrl     string      `json:"detailFileUrl"     orm:"detail_file_url"    description:""`                      //
	ExceptionDetail   string      `json:"exceptionDetail"   orm:"exception_detail"   description:""`                      //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"         description:""`                      //
	CompletedAt       *gtime.Time `json:"completedAt"       orm:"completed_at"       description:""`                      //
}
