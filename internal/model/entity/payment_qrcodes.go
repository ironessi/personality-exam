// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PaymentQrcodes is the golang structure for table payment_qrcodes.
type PaymentQrcodes struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`                       //
	TransactionNo string      `json:"transactionNo" orm:"transaction_no" description:"关联支付流水号"`                // 关联支付流水号
	QrType        int         `json:"qrType"        orm:"qr_type"        description:"类型：1-支付宝 2-微信"`          // 类型：1-支付宝 2-微信
	QrContent     string      `json:"qrContent"     orm:"qr_content"     description:"二维码内容"`                  // 二维码内容
	QrImageUrl    string      `json:"qrImageUrl"    orm:"qr_image_url"   description:"二维码图片URL"`               // 二维码图片URL
	Status        int         `json:"status"        orm:"status"         description:"0-有效 1-已扫描 2-已支付 3-已过期"` // 0-有效 1-已扫描 2-已支付 3-已过期
	ScannedAt     *gtime.Time `json:"scannedAt"     orm:"scanned_at"     description:"被扫描时间"`                  // 被扫描时间
	ScannedBy     string      `json:"scannedBy"     orm:"scanned_by"     description:"扫描者标识"`                  // 扫描者标识
	ExpireAt      int64       `json:"expireAt"      orm:"expire_at"      description:"过期时间"`                   // 过期时间
	ValidSeconds  uint        `json:"validSeconds"  orm:"valid_seconds"  description:""`                       //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""`                       //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`                       //
}
