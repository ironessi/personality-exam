package pay

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/guid"
)

type Dev struct{}

func NewDev() *Dev {
	return &Dev{}
}

type MockPaidInput struct {
	UserId            int64
	OrderNo           string
	Channel           string
	PaidTransactionId string
}

type MockPaidOutput struct {
	OrderNo           string
	TransactionNo     string
	PaidChannel       string
	PaidTransactionId string
}

// MockPaid 仅开发环境使用：把订单/流水直接置为“已支付”。
func (d *Dev) MockPaid(ctx context.Context, in MockPaidInput) (*MockPaidOutput, error) {
	var order entity.TestOrders
	if err := dao.TestOrders.Ctx(ctx).
		Where("order_no", in.OrderNo).
		Where("user_id", in.UserId).
		Scan(&order); err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, gerror.New("订单不存在")
	}
	if order.Status != 0 {
		return nil, gerror.New("仅支持待支付订单模拟支付")
	}

	paidTx := in.PaidTransactionId
	if paidTx == "" {
		paidTx = fmt.Sprintf("MOCK_%s", guid.S())
	}
	now := gtime.Now()

	var trans entity.PaymentTransactions
	_ = dao.PaymentTransactions.Ctx(ctx).
		Where("order_no", in.OrderNo).
		Where("status", 0).
		Order("created_at desc").
		Limit(1).
		Scan(&trans)

	// 若没有待支付流水，创建一条
	if trans.Id == 0 {
		transNo := generateTransNo(uint64(in.UserId))
		expireAt := gtime.Timestamp() + 300
		_, err := dao.PaymentTransactions.Ctx(ctx).Data(do.PaymentTransactions{
			TransactionNo: transNo,
			OrderNo:       in.OrderNo,
			UserId:        in.UserId,
			Amount:        order.PayableAmount,
			Channel:       in.Channel,
			Status:        0,
			ExpireAt:      expireAt,
		}).Insert()
		if err != nil {
			return nil, err
		}
		_ = dao.PaymentTransactions.Ctx(ctx).Where("transaction_no", transNo).Scan(&trans)
	}

	err := dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.PaymentTransactions.Ctx(ctx).
			WherePri(trans.Id).
			Data(do.PaymentTransactions{
				Status: 2,
				PaidAt: now,
				AlipayTradeNo: func() string {
					if in.Channel == "alipay" {
						return paidTx
					}
					return ""
				}(),
				WechatTransactionId: func() string {
					if in.Channel == "wechat" {
						return paidTx
					}
					return ""
				}(),
			}).Update()
		if err != nil {
			return err
		}

		_, err = dao.TestOrders.Ctx(ctx).
			Where("order_no", in.OrderNo).
			Where("status", 0).
			Data(do.TestOrders{
				Status:            1,
				PaidAt:            now,
				PaidChannel:       in.Channel,
				PaidTransactionId: paidTx,
			}).Update()
		return err
	})
	if err != nil {
		return nil, err
	}

	return &MockPaidOutput{
		OrderNo:           in.OrderNo,
		TransactionNo:     trans.TransactionNo,
		PaidChannel:       in.Channel,
		PaidTransactionId: paidTx,
	}, nil
}
