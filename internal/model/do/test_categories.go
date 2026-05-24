// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TestCategories is the golang structure of table test_categories for DAO operations like Where/Data.
type TestCategories struct {
	g.Meta      `orm:"table:test_categories, do:true"`
	Id          any         // 分类ID
	Name        any         // 分类名称
	Sort        any         // 排序优先级
	Status      any         // 状态(1-启用, 0-禁用)
	CreatedTime *gtime.Time // 创建时间
	UpdatedTime *gtime.Time // 更新时间
}
