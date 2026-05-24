// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TestTypes is the golang structure of table test_types for DAO operations like Where/Data.
type TestTypes struct {
	g.Meta      `orm:"table:test_types, do:true"`
	Id          any         // 类型ID，自增主键
	Name        any         // 测试名称
	Description any         // 测试描述
	TestTime    any         // 预计测试时间(分钟)
	Status      any         // 状态(0-禁用,1-启用)
	CreatedTime *gtime.Time // 创建时间
	UpdatedTime *gtime.Time // 更新时间
	Img         any         // 测试图标URL
	CategoryId  any         // 所属分类ID
	Price       any         // 查看报告所需支付金额（元）
}
