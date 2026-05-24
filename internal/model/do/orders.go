// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Orders is the golang structure of table orders for DAO operations like Where/Data.
type Orders struct {
	g.Meta       `orm:"table:orders, do:true"`
	Id           any         // 订单ID，自增主键
	UserId       any         // 所属用户ID
	TestTypeId   any         // 关联的问卷ID
	TestRecordId any         // 关联的测试记录ID
	OrderNo      any         // 商户订单号
	TradeNo      any         // 支付平台交易单号
	Amount       any         // 订单总金额（元）
	Status       any         // 0-待支付 1-已支付 2-已关闭
	PayType      any         // 支付渠道：alipay/wechat
	PayTime      *gtime.Time // 支付完成时间
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
