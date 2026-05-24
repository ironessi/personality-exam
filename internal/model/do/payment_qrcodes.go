// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentQrcodes is the golang structure of table payment_qrcodes for DAO operations like Where/Data.
type PaymentQrcodes struct {
	g.Meta        `orm:"table:payment_qrcodes, do:true"`
	Id            any         //
	TransactionNo any         // 关联支付流水号
	QrType        any         // 类型：1-支付宝 2-微信
	QrContent     any         // 二维码内容
	QrImageUrl    any         // 二维码图片URL
	Status        any         // 0-有效 1-已扫描 2-已支付 3-已过期
	ScannedAt     *gtime.Time // 被扫描时间
	ScannedBy     any         // 扫描者标识
	ExpireAt      any         // 过期时间
	ValidSeconds  any         //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
