package v1

import "github.com/gogf/gf/v2/frame/g"

// DevMockPaidReq 仅开发环境使用：模拟支付成功（无需真实回调）。
type DevMockPaidReq struct {
	g.Meta            `path:"/payment/dev/mock-paid" method:"post" summary:"(DEV)模拟支付成功" tags:"支付"`
	OrderNo           string `json:"orderNo" v:"required#订单号不能为空"`
	Channel           string `json:"channel" v:"required|in:alipay,wechat#支付渠道不能为空#支付渠道仅支持alipay或wechat"`
	PaidTransactionId string `json:"paidTransactionId" dc:"第三方交易号(可选)"`
}

type DevMockPaidRes struct {
	OrderNo           string `json:"orderNo"`
	TransactionNo     string `json:"transactionNo"`
	PaidChannel       string `json:"paidChannel"`
	PaidTransactionId string `json:"paidTransactionId"`
}
