// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Favorites is the golang structure for table favorites.
type Favorites struct {
	Id          int64       `json:"id"          orm:"id"           description:"收藏ID，自增主键"`  // 收藏ID，自增主键
	UserId      int64       `json:"userId"      orm:"user_id"      description:"关联用户ID"`     // 关联用户ID
	TestTypeId  int64       `json:"testTypeId"  orm:"test_type_id" description:"测试类型ID"`     // 测试类型ID
	ResultId    int64       `json:"resultId"    orm:"result_id"    description:"测试结果ID（可选）"` // 测试结果ID（可选）
	CreatedTime *gtime.Time `json:"createdTime" orm:"created_time" description:"收藏创建时间"`     // 收藏创建时间
	UpdatedTime *gtime.Time `json:"updatedTime" orm:"updated_time" description:"更新时间"`       // 更新时间
	DeletedTime *gtime.Time `json:"deletedTime" orm:"deleted_time" description:"软删除时间"`      // 软删除时间
}
