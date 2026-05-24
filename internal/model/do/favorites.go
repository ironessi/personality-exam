// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Favorites is the golang structure of table favorites for DAO operations like Where/Data.
type Favorites struct {
	g.Meta      `orm:"table:favorites, do:true"`
	Id          any         // 收藏ID，自增主键
	UserId      any         // 关联用户ID
	TestTypeId  any         // 测试类型ID
	ResultId    any         // 测试结果ID（可选）
	CreatedTime *gtime.Time // 收藏创建时间
	UpdatedTime *gtime.Time // 更新时间
	DeletedTime *gtime.Time // 软删除时间
}
