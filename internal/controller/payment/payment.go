package payment

import (
	"fmt"

	apiv1 "PersonalityTest/api/payment/v1"
	"PersonalityTest/internal/logic/pay"
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/shopspring/decimal"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) CreateOrder(ctx context.Context, req *apiv1.CreateOrderReq) (res *apiv1.CreateOrderRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	idempotencyKey := req.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = fmt.Sprintf("PAY_%d_%d_%s", userId, req.TestId, guid.S()[:6])
	}

	result, err := pay.NewOrders().CreateOrder(ctx, pay.CreateOrderInput{
		UserId:    userId,
		TestId:    uint(req.TestId),
		Channel:   req.Channel,
		ClientIp:  req.ClientIp,
		UserAgent: g.RequestFromCtx(ctx).Header.Get("User-Agent"),
		RecordId:  req.RecordId,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.CreateOrderRes{
		OrderId:        result.OrderId,
		OrderNo:        result.OrderNo,
		TestName:       result.TestName,
		OriginalPrice:  result.OriginalPrice,
		DiscountAmount: result.DiscountAmount,
		PayableAmount:  result.PayableAmount,
		Status:         result.Status,
		ExpireAt:       result.ExpireAt,
		CreatedAt:      result.CreatedAt,
	}, nil
}

func (c *Controller) DevMockPaid(ctx context.Context, req *apiv1.DevMockPaidReq) (res *apiv1.DevMockPaidRes, err error) {
	// 仅允许本地调用，防止误用到公网
	remoteIP := g.RequestFromCtx(ctx).GetClientIp()
	if !strings.HasPrefix(remoteIP, "127.0.0.1") && remoteIP != "::1" && remoteIP != "localhost" {
		return nil, gerror.New("仅允许本地调用该接口")
	}

	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	out, err := pay.NewDev().MockPaid(ctx, pay.MockPaidInput{
		UserId:            userId,
		OrderNo:           req.OrderNo,
		Channel:           req.Channel,
		PaidTransactionId: req.PaidTransactionId,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.DevMockPaidRes{
		OrderNo:           out.OrderNo,
		TransactionNo:     out.TransactionNo,
		PaidChannel:       out.PaidChannel,
		PaidTransactionId: out.PaidTransactionId,
	}, nil
}

func (c *Controller) GetOrderDetail(ctx context.Context, req *apiv1.GetOrderDetailReq) (res *apiv1.GetOrderDetailRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	order, err := pay.NewOrders().GetOrderByNo(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}

	if int64(order.UserId) != userId {
		return nil, gerror.New("订单不存在")
	}

	return &apiv1.GetOrderDetailRes{
		OrderId:           int64(order.Id),
		OrderNo:           order.OrderNo,
		UserId:            int64(order.UserId),
		TestId:            int(order.TestId),
		TestName:          order.TestName,
		OriginalPrice:     order.OriginalPrice,
		DiscountAmount:    order.DiscountAmount,
		PayableAmount:     order.PayableAmount,
		Status:            order.Status,
		PaidAt:            order.PaidAt,
		PaidChannel:       order.PaidChannel,
		PaidTransactionId: order.PaidTransactionId,
		ResultReportId:    int64(order.ResultReportId),
		ExpireAt:          order.ExpireAt,
		CreatedAt:         order.CreatedAt,
	}, nil
}

func (c *Controller) GetUserOrderList(ctx context.Context, req *apiv1.GetUserOrderListReq) (res *apiv1.GetUserOrderListRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	var status int = -1
	if req.Status != nil {
		status = *req.Status
	}

	result, err := pay.NewOrders().GetUserOrders(ctx, pay.GetUserOrdersInput{
		UserId:   userId,
		Status:   status,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	list := make([]apiv1.OrderListItem, len(result.List))
	for i, order := range result.List {
		list[i] = apiv1.OrderListItem{
			OrderNo:       order.OrderNo,
			TestName:      order.TestName,
			PayableAmount: order.PayableAmount,
			Status:        order.Status,
			PaidAt:        order.PaidAt,
			CreatedAt:     order.CreatedAt,
		}
	}

	return &apiv1.GetUserOrderListRes{
		List: list,
		Pagination: apiv1.Pagination{
			Total:      result.Total,
			Page:       result.Page,
			PageSize:   result.PageSize,
			TotalPages: result.TotalPages,
		},
	}, nil
}

func (c *Controller) CreateQrCode(ctx context.Context, req *apiv1.CreateQrCodeReq) (res *apiv1.CreateQrCodeRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	result, err := pay.NewQrcode().CreateQrcode(ctx, pay.CreateQrcodeInput{
		OrderNo: req.OrderNo,
		Channel: req.Channel,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.CreateQrCodeRes{
		QrCodeId:      result.QrCodeId,
		TransactionNo: result.TransactionNo,
		QrType:        result.QrType,
		QrContent:     result.QrContent,
		QrImageUrl:    result.QrImageUrl,
		Amount:        result.Amount,
		ExpireAt:      result.ExpireAt,
		ValidSeconds:  result.ValidSeconds,
	}, nil
}

func (c *Controller) GetQrCodeStatus(ctx context.Context, req *apiv1.GetQrCodeStatusReq) (res *apiv1.GetQrCodeStatusRes, err error) {
	result, err := pay.NewQrcode().GetQrcodeStatus(ctx, pay.GetQrcodeStatusInput{
		QrCodeId:      req.QrCodeId,
		TransactionNo: req.TransactionNo,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.GetQrCodeStatusRes{
		QrStatus:    result.QrStatus,
		PayStatus:   result.PayStatus,
		ScannedAt:   result.ScannedAt,
		PaidAt:      result.PaidAt,
		ExpireAt:    result.ExpireAt,
		FrontStatus: result.FrontStatus,
	}, nil
}

func (c *Controller) CancelOrder(ctx context.Context, req *apiv1.CancelOrderReq) (res *apiv1.CancelOrderRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	result, err := pay.NewOrders().CancelOrder(ctx, userId, req.OrderNo)
	if err != nil {
		return nil, err
	}

	return &apiv1.CancelOrderRes{
		OrderNo:     result.OrderNo,
		Status:      result.Status,
		CancelledAt: result.CancelledAt,
	}, nil
}

func (c *Controller) CheckExpire(ctx context.Context, req *apiv1.CheckExpireReq) (res *apiv1.CheckExpireRes, err error) {
	result, err := pay.NewOrders().CheckExpire(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}

	return &apiv1.CheckExpireRes{
		OrderNo:          result.OrderNo,
		IsExpired:        result.IsExpired,
		ExpireAt:         result.ExpireAt,
		RemainingSeconds: int64(result.RemainingSeconds),
	}, nil
}

func (c *Controller) CloseTransaction(ctx context.Context, req *apiv1.CloseTransactionReq) (res *apiv1.CloseTransactionRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	result, err := pay.NewOrders().CloseTransaction(ctx, pay.CloseTransactionInput{
		TransactionNo: req.TransactionNo,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.CloseTransactionRes{
		TransactionNo: result.TransactionNo,
		Status:        result.Status,
		ClosedAt:      result.ClosedAt,
	}, nil
}

func (c *Controller) CreateRefund(ctx context.Context, req *apiv1.CreateRefundReq) (res *apiv1.CreateRefundRes, err error) {
	userId := g.RequestFromCtx(ctx).GetCtxVar("uid").Int64()
	if userId == 0 {
		return nil, gerror.New("请先登录")
	}

	result, err := pay.NewRefunds().ApplyRefund(ctx, pay.ApplyRefundInput{
		UserId:  userId,
		OrderNo: req.OrderNo,
		Reason:  req.Reason,
	})
	if err != nil {
		return nil, err
	}

	return &apiv1.CreateRefundRes{
		RefundNo:     result.RefundNo,
		OrderNo:      result.OrderNo,
		RefundAmount: result.RefundAmount,
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
	}, nil
}

func (c *Controller) GetRefundDetail(ctx context.Context, req *apiv1.GetRefundDetailReq) (res *apiv1.GetRefundDetailRes, err error) {
	result, err := pay.NewRefunds().GetRefundByNo(ctx, req.RefundNo)
	if err != nil {
		return nil, err
	}

	return &apiv1.GetRefundDetailRes{
		RefundNo:        result.RefundNo,
		OrderNo:         result.OrderNo,
		RefundAmount:    result.RefundAmount,
		Reason:          result.Reason,
		Status:          result.Status,
		ChannelRefundId: result.ChannelRefundId,
		CompletedAt:     result.CompletedAt,
		CreatedAt:       result.CreatedAt,
	}, nil
}

func (c *Controller) GetChannels(ctx context.Context, req *apiv1.GetChannelsReq) (res *apiv1.GetChannelsRes, err error) {
	result, err := pay.NewChannels().GetChannels(ctx)
	if err != nil {
		return nil, err
	}

	channels := make([]apiv1.ChannelInfo, len(result.Channels))
	for i, ch := range result.Channels {
		channels[i] = apiv1.ChannelInfo{
			Channel:     ch.Channel,
			Name:        ch.Name,
			Icon:        ch.Icon,
			Enabled:     ch.Enabled,
			MinAmount:   ch.MinAmount,
			MaxAmount:   ch.MaxAmount,
			Description: ch.Description,
		}
	}

	return &apiv1.GetChannelsRes{
		Channels: channels,
	}, nil
}

func (c *Controller) GetTestPrice(ctx context.Context, req *apiv1.GetTestPriceReq) (res *apiv1.GetTestPriceRes, err error) {
	result, err := pay.NewPrice().GetTestPrice(ctx, pay.GetTestPriceInput{TestId: int64(req.TestId)})
	if err != nil {
		return nil, err
	}

	return &apiv1.GetTestPriceRes{
		TestId:        result.TestId,
		TestName:      result.TestName,
		OriginalPrice: result.OriginalPrice,
		DiscountPrice: result.DiscountPrice,
		DiscountRate:  result.DiscountRate,
		PromotionText: result.PromotionText,
	}, nil
}

func (c *Controller) AlipayNotify(r *ghttp.Request) {
	ctx := r.Context()
	if err := pay.New().HandleAlipayNotify(ctx); err != nil {
		g.Log().Errorf(ctx, "支付宝回调处理失败: %v", err)
		r.Response.WriteStatus(400, "fail")
		return
	}
	r.Response.Write("success")
}

func (c *Controller) WechatNotify(r *ghttp.Request) {
	ctx := r.Context()
	if err := pay.New().HandleWechatNotify(ctx); err != nil {
		g.Log().Errorf(ctx, "微信支付回调处理失败: %v", err)
		r.Response.WriteStatus(400, "fail")
		return
	}
	r.Response.Write("success")
}

func (c *Controller) QueryOrder(r *ghttp.Request) {
	ctx := r.Context()
	orderNo := r.Get("orderNo").String()
	order, err := pay.New().QueryOrder(ctx, orderNo)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 1, "message": err.Error()})
		return
	}
	r.Response.WriteJson(g.Map{
		"code":    0,
		"message": "",
		"data": g.Map{
			"orderNo": order.OrderNo,
			"status":  order.Status,
			"amount":  decimal.NewFromFloat(order.PayableAmount),
		},
	})
}

func (c *Controller) RefundOrder(r *ghttp.Request) {
	ctx := r.Context()
	userId := r.GetCtxVar("uid").Int64()
	if userId == 0 {
		r.Response.WriteJsonExit(g.Map{"code": 1, "message": "请先登录"})
		return
	}

	orderNo := r.Get("orderNo").String()
	reason := r.Get("reason").String()
	result, err := pay.NewRefunds().ApplyRefund(ctx, pay.ApplyRefundInput{
		UserId:  userId,
		OrderNo: orderNo,
		Reason:  reason,
	})
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 1, "message": err.Error()})
		return
	}
	r.Response.WriteJson(g.Map{
		"code":    0,
		"message": "",
		"data": g.Map{
			"refundNo": result.RefundNo,
		},
	})
}

func RegisterRoutes(group *ghttp.RouterGroup) {
	ctrl := New()

	group.Group("/", func(group *ghttp.RouterGroup) {
		// 绑定通过结构体 g.Meta 定义的严格路由分组
		group.Bind(ctrl)

		// 绑定没有使用 g.Meta 定义的普通非严格路由（保持原来地址）
		group.GET("/payment/query", ctrl.QueryOrder)
		// 兼容旧退款路由，但内部统一走新的退款申请逻辑。
		group.POST("/payment/refund", ctrl.RefundOrder)
	})
}
