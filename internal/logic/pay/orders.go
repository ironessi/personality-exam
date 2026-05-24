package pay

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type Orders struct{}

func NewOrders() *Orders {
	return &Orders{}
}

type CreateOrderInput struct {
	UserId    int64
	TestId    uint
	Channel   string
	ClientIp  string
	UserAgent string
	RecordId  int64
}

type CreateOrderOutput struct {
	OrderId        int64       `json:"orderId"`
	OrderNo        string      `json:"orderNo"`
	TestName       string      `json:"testName"`
	OriginalPrice  float64     `json:"originalPrice"`
	DiscountAmount float64     `json:"discountAmount"`
	PayableAmount  float64     `json:"payableAmount"`
	Status         int         `json:"status"`
	ExpireAt       int64       `json:"expireAt"`
	CreatedAt      *gtime.Time `json:"createdAt"`
}

func (o *Orders) CreateOrder(ctx context.Context, in CreateOrderInput) (*CreateOrderOutput, error) {
	var testType entity.TestTypes
	err := dao.TestTypes.Ctx(ctx).WherePri(in.TestId).Scan(&testType)
	if err != nil {
		return nil, err
	}
	if testType.Id == 0 {
		return nil, gerror.New("测试类型不存在")
	}

	var existOrder entity.TestOrders
	_ = dao.TestOrders.Ctx(ctx).
		Where("user_id", in.UserId).
		Where("test_id", in.TestId).
		Where("status", 1).
		Scan(&existOrder)
	if existOrder.Id > 0 {
		return nil, gerror.New("您已支付过该测试")
	}

	orderNo := generateOrderNo(uint64(in.UserId))
	expireAt := gtime.Timestamp() + 1800 // 30分钟

	err = dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.TestOrders.Ctx(ctx).Data(do.TestOrders{
			OrderNo:        orderNo,
			UserId:         in.UserId,
			TestId:         in.TestId,
			TestName:       testType.Name,
			OriginalPrice:  testType.Price,
			DiscountAmount: 0,
			PayableAmount:  testType.Price,
			Status:         0,
			ResultReportId: uint64(in.RecordId),
			ExpireAt:       expireAt,
			ClientIp:       in.ClientIp,
			UserAgent:      in.UserAgent,
		}).Insert()
		return err
	})
	if err != nil {
		return nil, err
	}

	var order entity.TestOrders
	dao.TestOrders.Ctx(ctx).Where("order_no", orderNo).Scan(&order)

	return &CreateOrderOutput{
		OrderId:        int64(order.Id),
		OrderNo:        order.OrderNo,
		TestName:       order.TestName,
		OriginalPrice:  order.OriginalPrice,
		DiscountAmount: order.DiscountAmount,
		PayableAmount:  order.PayableAmount,
		Status:         order.Status,
		ExpireAt:       order.ExpireAt,
		CreatedAt:      order.CreatedAt,
	}, nil
}

func (o *Orders) GetOrderByNo(ctx context.Context, orderNo string) (*entity.TestOrders, error) {
	var order entity.TestOrders
	err := dao.TestOrders.Ctx(ctx).Where("order_no", orderNo).Scan(&order)
	if err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, gerror.New("订单不存在")
	}
	return &order, nil
}

type GetUserOrdersInput struct {
	UserId   int64
	Status   int
	Page     int
	PageSize int
}

type GetUserOrdersOutput struct {
	List       []*entity.TestOrders `json:"list"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"pageSize"`
	TotalPages int                  `json:"totalPages"`
}

func (o *Orders) GetUserOrders(ctx context.Context, in GetUserOrdersInput) (*GetUserOrdersOutput, error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 || in.PageSize > 100 {
		in.PageSize = 10
	}

	m := dao.TestOrders.Ctx(ctx).Where("user_id", in.UserId).Where("is_deleted", 0)
	if in.Status >= 0 {
		m = m.Where("status", in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var list []*entity.TestOrders
	err = m.Order("created_at desc").Page(in.Page, in.PageSize).Scan(&list)
	if err != nil {
		return nil, err
	}

	totalPages := total / in.PageSize
	if total%in.PageSize > 0 {
		totalPages++
	}

	return &GetUserOrdersOutput{
		List:       list,
		Total:      total,
		Page:       in.Page,
		PageSize:   in.PageSize,
		TotalPages: totalPages,
	}, nil
}

type CancelOrderOutput struct {
	OrderNo     string      `json:"orderNo"`
	Status      int         `json:"status"`
	CancelledAt *gtime.Time `json:"cancelledAt"`
}

func (o *Orders) CancelOrder(ctx context.Context, userId int64, orderNo string) (*CancelOrderOutput, error) {
	var order entity.TestOrders
	err := dao.TestOrders.Ctx(ctx).
		Where("order_no", orderNo).
		Where("user_id", userId).
		Scan(&order)
	if err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, gerror.New("订单不存在")
	}
	if order.Status != 0 {
		return nil, gerror.New("只能取消待支付的订单")
	}

	cancelledAt := gtime.Now()
	_, err = dao.TestOrders.Ctx(ctx).
		Where("order_no", orderNo).
		Data(do.TestOrders{
			Status: 5,
		}).Update()
	if err != nil {
		return nil, err
	}

	return &CancelOrderOutput{
		OrderNo:     orderNo,
		Status:      5,
		CancelledAt: cancelledAt,
	}, nil
}

type CheckExpireOutput struct {
	OrderNo          string `json:"orderNo"`
	IsExpired        bool   `json:"isExpired"`
	ExpireAt         int64  `json:"expireAt"`
	RemainingSeconds int    `json:"remainingSeconds"`
}

func (o *Orders) CheckExpire(ctx context.Context, orderNo string) (*CheckExpireOutput, error) {
	var order entity.TestOrders
	err := dao.TestOrders.Ctx(ctx).Where("order_no", orderNo).Scan(&order)
	if err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, gerror.New("订单不存在")
	}

	remaining := int(order.ExpireAt - gtime.Timestamp())

	return &CheckExpireOutput{
		OrderNo:          orderNo,
		IsExpired:        remaining <= 0,
		ExpireAt:         order.ExpireAt,
		RemainingSeconds: remaining,
	}, nil
}

type CloseTransactionInput struct {
	TransactionNo string `json:"transactionNo"`
}

type CloseTransactionOutput struct {
	TransactionNo string      `json:"transactionNo"`
	Status        int         `json:"status"`
	ClosedAt      *gtime.Time `json:"closedAt"`
}

func (o *Orders) CloseTransaction(ctx context.Context, in CloseTransactionInput) (*CloseTransactionOutput, error) {
	var trans entity.PaymentTransactions
	err := dao.PaymentTransactions.Ctx(ctx).Where("transaction_no", in.TransactionNo).Scan(&trans)
	if err != nil {
		return nil, err
	}
	if trans.Id == 0 {
		return nil, gerror.New("支付流水不存在")
	}

	// 只有待支付(0)或支付中(1)的流水可以关闭
	if trans.Status != 0 && trans.Status != 1 {
		return nil, gerror.New("当前状态不允许关闭")
	}

	now := gtime.Now()
	_, err = dao.PaymentTransactions.Ctx(ctx).Where("id", trans.Id).Data(map[string]interface{}{
		"status":     4, // 已关闭
		"closed_at":  now,
		"updated_at": now,
	}).Update()
	if err != nil {
		return nil, err
	}

	return &CloseTransactionOutput{
		TransactionNo: in.TransactionNo,
		Status:        4,
		ClosedAt:      now,
	}, nil
}
