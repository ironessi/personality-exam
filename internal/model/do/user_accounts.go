// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type UserAccounts struct {
	g.Meta         `orm:"table:user_accounts, do:true"`
	Id             any         //
	UserId         any         // 用户ID
	Balance        any         // 余额（元）
	FrozenAmount   any         // 冻结金额
	TotalRecharge  any         // 累计充值
	TotalConsume   any         // 累计消费
	Version        any         // 乐观锁版本号
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
