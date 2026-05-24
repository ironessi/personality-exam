// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TestTypes is the golang structure for table test_types.
type TestTypes struct {
	Id          int64       `json:"id"          orm:"id"           description:"类型ID，自增主键"`     // 类型ID，自增主键
	Name        string      `json:"name"        orm:"name"         description:"测试名称"`          // 测试名称
	Description string      `json:"description" orm:"description"  description:"测试描述"`          // 测试描述
	TestTime    int         `json:"testTime"    orm:"test_time"    description:"预计测试时间(分钟)"`    // 预计测试时间(分钟)
	Status      int         `json:"status"      orm:"status"       description:"状态(0-禁用,1-启用)"` // 状态(0-禁用,1-启用)
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"创建时间"`          // 创建时间
	UpdatedTime *gtime.Time `json:"updatedTime" orm:"updated_time" description:"更新时间"`          // 更新时间
	Img         string      `json:"img"         orm:"img"          description:"测试图标URL"`       // 测试图标URL
	CategoryId  int64       `json:"categoryId"  orm:"category_id"  description:"所属分类ID"`        // 所属分类ID
	Price       float64     `json:"price"       orm:"price"        description:"查看报告所需支付金额（元）"` // 查看报告所需支付金额（元）
}
