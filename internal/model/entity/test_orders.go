// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TestOrders is the golang structure for table test_orders.
type TestOrders struct {
	Id                uint64      `json:"id"                orm:"id"                  description:"订单ID"`                                   // 订单ID
	OrderNo           string      `json:"orderNo"           orm:"order_no"            description:"订单编号（业务唯一）"`                             // 订单编号（业务唯一）
	UserId            uint64      `json:"userId"            orm:"user_id"             description:"用户ID"`                                   // 用户ID
	TestId            uint        `json:"testId"            orm:"test_id"             description:"测试类型ID"`                                 // 测试类型ID
	TestName          string      `json:"testName"          orm:"test_name"           description:"测试名称（冗余存储）"`                             // 测试名称（冗余存储）
	OriginalPrice     float64     `json:"originalPrice"     orm:"original_price"      description:"原价"`                                     // 原价
	DiscountAmount    float64     `json:"discountAmount"    orm:"discount_amount"     description:"优惠金额"`                                   // 优惠金额
	PayableAmount     float64     `json:"payableAmount"     orm:"payable_amount"      description:"应付金额（实际支付）"`                             // 应付金额（实际支付）
	Status            int         `json:"status"            orm:"status"              description:"状态：0-待支付 1-已支付 2-处理中 3-已完成 4-已退款 5-已关闭"` // 状态：0-待支付 1-已支付 2-处理中 3-已完成 4-已退款 5-已关闭
	PaidAt            *gtime.Time `json:"paidAt"            orm:"paid_at"             description:"支付完成时间"`                                 // 支付完成时间
	PaidChannel       string      `json:"paidChannel"       orm:"paid_channel"        description:"支付渠道：alipay/wechat"`                     // 支付渠道：alipay/wechat
	PaidTransactionId string      `json:"paidTransactionId" orm:"paid_transaction_id" description:"第三方支付流水号"`                               // 第三方支付流水号
	ResultReportId    uint64      `json:"resultReportId"    orm:"result_report_id"    description:"生成的报告ID"`                                // 生成的报告ID
	ExpireAt          int64       `json:"expireAt"          orm:"expire_at"           description:"订单过期时间"`                                 // 订单过期时间
	ClientIp          string      `json:"clientIp"          orm:"client_ip"           description:"用户下单IP"`                                 // 用户下单IP
	UserAgent         string      `json:"userAgent"         orm:"user_agent"          description:"用户代理"`                                   // 用户代理
	IsDeleted         int         `json:"isDeleted"         orm:"is_deleted"          description:""`                                       //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          description:""`                                       //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          description:""`                                       //
}
