// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// Orders is the golang structure for table orders.
type Orders struct {
	Id           int64           `json:"id"           orm:"id"             description:"订单ID，自增主键"`
	UserId       int64           `json:"userId"       orm:"user_id"         description:"所属用户ID"`
	TestTypeId   int64           `json:"testTypeId"   orm:"test_type_id"    description:"关联的问卷ID"`
	TestRecordId int64           `json:"testRecordId" orm:"test_record_id"  description:"关联的测试记录ID"`
	OrderNo      string          `json:"orderNo"      orm:"order_no"        description:"商户订单号（唯一）"`
	TradeNo      string          `json:"tradeNo"      orm:"trade_no"        description:"支付平台交易单号"`
	Amount       decimal.Decimal `json:"amount"       orm:"amount"          description:"订单总金额（元）"`
	Status       int             `json:"status"       orm:"status"          description:"0-待支付 1-已支付 2-已关闭"`
	PayType      string          `json:"payType"      orm:"pay_type"        description:"支付渠道：alipay/wechat"`
	PayTime      *gtime.Time     `json:"payTime"      orm:"pay_time"        description:"支付完成时间"`
	CreatedAt    *gtime.Time     `json:"createdAt"    orm:"created_at"      description:"创建时间"`
	UpdatedAt    *gtime.Time     `json:"updatedAt"    orm:"updated_at"      description:"更新时间"`
}
