// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ReconciliationRecords is the golang structure of table reconciliation_records for DAO operations like Where/Data.
type ReconciliationRecords struct {
	g.Meta            `orm:"table:reconciliation_records, do:true"`
	Id                any         //
	ReconcileDate     *gtime.Time // 对账日期
	Channel           any         // 渠道
	TotalTransactions any         //
	TotalAmount       any         //
	MatchedCount      any         //
	UnmatchedCount    any         //
	Status            any         // 0-待对账 1-对账中 2-完成 3-异常
	DetailFileUrl     any         //
	ExceptionDetail   any         //
	CreatedAt         *gtime.Time //
	CompletedAt       *gtime.Time //
}
