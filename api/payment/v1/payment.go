package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type CreateOrderReq struct {
	g.Meta         `path:"/payment/orders" method:"post" summary:"创建测试订单" tags:"支付"`
	TestId         int64  `json:"testId" v:"required#测试ID不能为空"`
	Channel        string `json:"channel" v:"required|in:alipay,wechat#支付渠道不能为空#支付渠道仅支持alipay或wechat"`
	ClientIp       string `json:"clientIp"`
	IdempotencyKey string `json:"idempotencyKey"`
	RecordId       int64  `json:"recordId"`
}

type CreateOrderRes struct {
	OrderId        int64        `json:"orderId"`
	OrderNo        string       `json:"orderNo"`
	TestName       string       `json:"testName"`
	OriginalPrice  float64      `json:"originalPrice"`
	DiscountAmount float64      `json:"discountAmount"`
	PayableAmount  float64      `json:"payableAmount"`
	Status         int          `json:"status"`
	ExpireAt       int64        `json:"expireAt"`
	CreatedAt      *gtime.Time  `json:"createdAt"`
}

type GetOrderDetailReq struct {
	g.Meta  `path:"/payment/orders/{orderNo}" method:"get" summary:"获取订单详情" tags:"支付"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空" dc:"订单编号"`
}

type GetOrderDetailRes struct {
	OrderId             int64        `json:"orderId"`
	OrderNo             string       `json:"orderNo"`
	UserId              int64        `json:"userId"`
	TestId              int          `json:"testId"`
	TestName            string       `json:"testName"`
	OriginalPrice       float64      `json:"originalPrice"`
	DiscountAmount      float64      `json:"discountAmount"`
	PayableAmount       float64      `json:"payableAmount"`
	Status              int          `json:"status"`
	PaidAt              *gtime.Time  `json:"paidAt"`
	PaidChannel         string       `json:"paidChannel"`
	PaidTransactionId   string       `json:"paidTransactionId"`
	ResultReportId      int64        `json:"resultReportId"`
	ExpireAt            int64        `json:"expireAt"`
	CreatedAt           *gtime.Time  `json:"createdAt"`
}

type GetUserOrderListReq struct {
	g.Meta    `path:"/payment/orders/list" method:"get" summary:"获取用户订单列表" tags:"支付"`
	Page      int    `json:"page" d:"1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" dc:"每页数量"`
	Status    *int   `json:"status" dc:"订单状态：0-待支付 1-已支付 2-处理中 3-已完成 4-已退款 5-已关闭"`
	StartTime string `json:"startTime" dc:"开始时间（ISO 8601）"`
	EndTime   string `json:"endTime" dc:"结束时间（ISO 8601）"`
}

type OrderListItem struct {
	OrderNo        string       `json:"orderNo"`
	TestName       string       `json:"testName"`
	PayableAmount  float64      `json:"payableAmount"`
	Status         int          `json:"status"`
	PaidAt        *gtime.Time  `json:"paidAt"`
	CreatedAt     *gtime.Time  `json:"createdAt"`
}

type Pagination struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}

type GetUserOrderListRes struct {
	List       []OrderListItem `json:"list"`
	Pagination Pagination       `json:"pagination"`
}

type CreateQrCodeReq struct {
	g.Meta  `path:"/payment/qrcode" method:"post" summary:"创建支付二维码" tags:"支付"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空"`
	Channel string `json:"channel" v:"required|in:alipay,wechat#支付渠道不能为空#支付渠道仅支持alipay或wechat"`
}

type CreateQrCodeRes struct {
	QrCodeId      string       `json:"qrCodeId"`
	TransactionNo string       `json:"transactionNo"`
	QrType        int          `json:"qrType"`
	QrContent     string       `json:"qrContent"`
	QrImageUrl    string       `json:"qrImageUrl"`
	Amount        float64      `json:"amount"`
	ExpireAt      int64        `json:"expireAt"`
	ValidSeconds  int          `json:"validSeconds"`
}

type GetQrCodeStatusReq struct {
	g.Meta         `path:"/payment/qrcode/status" method:"get" summary:"查询二维码支付状态" tags:"支付"`
	QrCodeId       string `json:"qrCodeId" v:"required#二维码ID不能为空"`
	TransactionNo  string `json:"transactionNo" v:"required#支付流水号不能为空"`
}

type GetQrCodeStatusRes struct {
	QrStatus       int          `json:"qrStatus"`
	PayStatus      int          `json:"payStatus"`
	ScannedAt      *gtime.Time  `json:"scannedAt"`
	PaidAt         *gtime.Time  `json:"paidAt"`
	ExpireAt       int64        `json:"expireAt"`
	FrontStatus    string       `json:"frontStatus"`
	ResultReportId int64        `json:"resultReportId"`
}

type CancelOrderReq struct {
	g.Meta  `path:"/payment/orders/{orderNo}/cancel" method:"post" summary:"取消订单" tags:"支付"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空" dc:"订单编号"`
}

type CancelOrderRes struct {
	OrderNo     string       `json:"orderNo"`
	Status      int          `json:"status"`
	CancelledAt *gtime.Time `json:"cancelledAt"`
}

type CheckExpireReq struct {
	g.Meta  `path:"/payment/orders/{orderNo}/check-expire" method:"get" summary:"检查订单过期" tags:"支付"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空" dc:"订单编号"`
}

type CheckExpireRes struct {
	OrderNo          string       `json:"orderNo"`
	IsExpired        bool         `json:"isExpired"`
	ExpireAt         int64        `json:"expireAt"`
	RemainingSeconds int64        `json:"remainingSeconds"`
}

type CreateRefundReq struct {
	g.Meta  `path:"/payment/refunds" method:"post" summary:"申请退款" tags:"支付"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空"`
	Reason  string `json:"reason" v:"required#退款原因不能为空"`
}

type CreateRefundRes struct {
	RefundNo    string       `json:"refundNo"`
	OrderNo     string       `json:"orderNo"`
	RefundAmount float64     `json:"refundAmount"`
	Status      int          `json:"status"`
	CreatedAt   *gtime.Time `json:"createdAt"`
}

type GetRefundDetailReq struct {
	g.Meta   `path:"/payment/refunds/{refundNo}" method:"get" summary:"查询退款记录" tags:"支付"`
	RefundNo string `json:"refundNo" v:"required#退款单号不能为空" dc:"退款单号"`
}

type GetRefundDetailRes struct {
	RefundNo        string       `json:"refundNo"`
	OrderNo         string       `json:"orderNo"`
	RefundAmount    float64     `json:"refundAmount"`
	Reason          string       `json:"reason"`
	Status          int         `json:"status"`
	ChannelRefundId string       `json:"channelRefundId"`
	CompletedAt    *gtime.Time  `json:"completedAt"`
	CreatedAt       *gtime.Time `json:"createdAt"`
}

type CloseTransactionReq struct {
	g.Meta        `path:"/payment/transactions/{transactionNo}/close" method:"post" summary:"关闭支付流水" tags:"支付"`
	TransactionNo string `json:"transactionNo" v:"required#支付流水号不能为空" dc:"支付流水号"`
}

type CloseTransactionRes struct {
	TransactionNo string      `json:"transactionNo"`
	Status        int         `json:"status"`
	ClosedAt      *gtime.Time `json:"closedAt"`
}

type GetChannelsReq struct {
	g.Meta `path:"/payment/channels" method:"get" summary:"获取支付渠道配置" tags:"支付"`
}

type ChannelInfo struct {
	Channel     string  `json:"channel"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	Enabled     bool    `json:"enabled"`
	MinAmount   float64 `json:"minAmount"`
	MaxAmount   float64 `json:"maxAmount"`
	Description string  `json:"description"`
}

type GetChannelsRes struct {
	Channels []ChannelInfo `json:"channels"`
}

type GetTestPriceReq struct {
	g.Meta `path:"/payment/tests/{testId}/price" method:"get" summary:"获取测试价格信息" tags:"支付"`
	TestId int64 `json:"testId" v:"required#测试ID不能为空" dc:"测试类型ID"`
}

type GetTestPriceRes struct {
	TestId         int64   `json:"testId"`
	TestName       string  `json:"testName"`
	OriginalPrice  float64 `json:"originalPrice"`
	DiscountPrice  float64 `json:"discountPrice"`
	DiscountRate   float64 `json:"discountRate"`
	PromotionText  string  `json:"promotionText"`
}
