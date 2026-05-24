// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TestOrders is the golang structure of table test_orders for DAO operations like Where/Data.
type TestOrders struct {
	g.Meta            `orm:"table:test_orders, do:true"`
	Id                any         // 订单ID
	OrderNo           any         // 订单编号（业务唯一）
	UserId            any         // 用户ID
	TestId            any         // 测试类型ID
	TestName          any         // 测试名称（冗余存储）
	OriginalPrice     any         // 原价
	DiscountAmount    any         // 优惠金额
	PayableAmount     any         // 应付金额（实际支付）
	Status            any         // 状态：0-待支付 1-已支付 2-处理中 3-已完成 4-已退款 5-已关闭
	PaidAt            *gtime.Time // 支付完成时间
	PaidChannel       any         // 支付渠道：alipay/wechat
	PaidTransactionId any         // 第三方支付流水号
	ResultReportId    any         // 生成的报告ID
	ExpireAt          any         // 订单过期时间
	ClientIp          any         // 用户下单IP
	UserAgent         any         // 用户代理
	IsDeleted         any         //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
