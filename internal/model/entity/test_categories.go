// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TestCategories is the golang structure for table test_categories.
type TestCategories struct {
	Id          int64       `json:"id"          orm:"id"           description:"分类ID"`           // 分类ID
	Name        string      `json:"name"        orm:"name"         description:"分类名称"`           // 分类名称
	Sort        int         `json:"sort"        orm:"sort"         description:"排序优先级"`          // 排序优先级
	Status      int         `json:"status"      orm:"status"       description:"状态(1-启用, 0-禁用)"` // 状态(1-启用, 0-禁用)
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"创建时间"`           // 创建时间
	UpdatedTime *gtime.Time `json:"updatedTime" orm:"updated_time" description:"更新时间"`           // 更新时间
}
