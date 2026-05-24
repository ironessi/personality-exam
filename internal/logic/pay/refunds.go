package pay

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"
	"fmt"
	"time"

	"github.com/go-pay/gopay"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

type Refunds struct{}

func NewRefunds() *Refunds {
	return &Refunds{}
}

type ApplyRefundInput struct {
	UserId  int64
	OrderNo string
	Reason  string
}

type ApplyRefundOutput struct {
	RefundNo     string      `json:"refundNo"`
	OrderNo      string      `json:"orderNo"`
	RefundAmount float64     `json:"refundAmount"`
	Status       int         `json:"status"`
	CreatedAt    *gtime.Time `json:"createdAt"`
}

func (r *Refunds) ApplyRefund(ctx context.Context, in ApplyRefundInput) (*ApplyRefundOutput, error) {
	refundNo := generateRefundNo(in.UserId)
	createdAt := gtime.Now()
	output := &ApplyRefundOutput{
		RefundNo:  refundNo,
		OrderNo:   in.OrderNo,
		Status:    0,
		CreatedAt: createdAt,
	}

	err := dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var order entity.TestOrders
		err := dao.TestOrders.Ctx(ctx).
			Where("order_no", in.OrderNo).
			Where("user_id", in.UserId).
			LockUpdate().
			Scan(&order)
		if err != nil {
			return err
		}
		if order.Id == 0 {
			return gerror.New("订单不存在")
		}
		if order.Status != 1 && order.Status != 3 {
			return gerror.New("当前订单状态不支持退款")
		}

		var trans entity.PaymentTransactions
		err = dao.PaymentTransactions.Ctx(ctx).
			Where("order_no", in.OrderNo).
			Where("status", 2).
			LockUpdate().
			Scan(&trans)
		if err != nil {
			return err
		}
		if trans.Id == 0 {
			return gerror.New("未找到成功的支付流水")
		}

		var existRefund entity.RefundRecords
		err = dao.RefundRecords.Ctx(ctx).
			Where("order_no", in.OrderNo).
			Where("status IN (0,1,2)").
			LockUpdate().
			OrderDesc("id").
			Limit(1).
			Scan(&existRefund)
		if err != nil {
			return err
		}
		if existRefund.Id > 0 {
			return gerror.New("该订单已有退款申请")
		}

		_, err = dao.RefundRecords.Ctx(ctx).Data(do.RefundRecords{
			RefundNo:      refundNo,
			TransactionNo: trans.TransactionNo,
			OrderNo:       in.OrderNo,
			RefundAmount:  trans.Amount,
			Reason:        in.Reason,
			Status:        0,
		}).Insert()
		if err != nil {
			return err
		}

		output.RefundAmount = trans.Amount
		return nil
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (r *Refunds) ProcessPendingRefunds(ctx context.Context, limit int) {
	if limit <= 0 {
		limit = 10
	}

	var pending []entity.RefundRecords
	err := dao.RefundRecords.Ctx(ctx).
		Where("status", 0).
		Order("id asc").
		Limit(limit).
		Scan(&pending)
	if err != nil {
		g.Log().Errorf(ctx, "扫描待处理退款失败: %v", err)
		return
	}

	for _, refund := range pending {
		if err := r.processRefund(ctx, int64(refund.Id)); err != nil {
			g.Log().Errorf(ctx, "处理退款失败 [refundNo=%s]: %v", refund.RefundNo, err)
		}
	}
}

func (r *Refunds) processRefund(ctx context.Context, refundId int64) error {
	var refund entity.RefundRecords
	var trans entity.PaymentTransactions
	var order entity.TestOrders

	err := dao.RefundRecords.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		affected, err := dao.RefundRecords.Ctx(ctx).
			WherePri(refundId).
			Where("status", 0).
			Data(do.RefundRecords{
				Status: 1,
			}).Update()
		if err != nil {
			return err
		}

		rowsAffected, _ := affected.RowsAffected()
		if rowsAffected == 0 {
			return nil
		}

		err = dao.RefundRecords.Ctx(ctx).
			WherePri(refundId).
			Scan(&refund)
		if err != nil {
			return err
		}
		if refund.Id == 0 {
			return gerror.New("退款记录不存在")
		}

		err = dao.PaymentTransactions.Ctx(ctx).
			Where("transaction_no", refund.TransactionNo).
			LockUpdate().
			Scan(&trans)
		if err != nil {
			return err
		}
		if trans.Id == 0 {
			return gerror.New("未找到对应支付流水")
		}

		err = dao.TestOrders.Ctx(ctx).
			Where("order_no", refund.OrderNo).
			LockUpdate().
			Scan(&order)
		if err != nil {
			return err
		}
		if order.Id == 0 {
			return gerror.New("未找到对应订单")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if refund.Id == 0 {
		return nil
	}

	var refundErr error
	if trans.Channel == "alipay" {
		refundErr = r.processAlipayRefund(ctx, &trans, decimal.NewFromFloat(trans.Amount), refund.RefundNo, refund.Reason)
	} else if trans.Channel == "wechat" {
		refundErr = gerror.New("微信退款功能正在开发中")
	} else {
		refundErr = gerror.New("不支持的支付渠道")
	}

	if refundErr != nil {
		_, updateErr := dao.RefundRecords.Ctx(ctx).WherePri(refund.Id).Data(g.Map{
			"status":           3,
			"channel_response": refundErr.Error(),
		}).Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, refundErr.Error())
		}
		return refundErr
	}

	return dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.RefundRecords.Ctx(ctx).WherePri(refund.Id).Data(g.Map{
			"status":       2,
			"completed_at": gtime.Now(),
		}).Update()
		if err != nil {
			return err
		}

		_, err = dao.TestOrders.Ctx(ctx).
			Where("order_no", order.OrderNo).
			Data(g.Map{"status": 4}).
			Update()
		if err != nil {
			return err
		}

		_, err = dao.PaymentTransactions.Ctx(ctx).
			Where("order_no", order.OrderNo).
			Where("status", 2).
			Data(g.Map{"status": 5}).
			Update()
		return err
	})
}

func (r *Refunds) processAlipayRefund(ctx context.Context, trans *entity.PaymentTransactions, amount decimal.Decimal, refundNo, reason string) error {
	client, _, err := NewAlipayClient(ctx)
	if err != nil {
		return err
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", trans.OrderNo).
		Set("refund_amount", amount.StringFixed(2)).
		Set("refund_reason", reason).
		Set("out_request_no", refundNo)

	res, err := client.TradeRefund(ctx, bm)
	if err != nil {
		return err
	}
	if res.Response.Code != "10000" {
		return gerror.New(res.Response.Msg + " " + res.Response.SubMsg)
	}

	return nil
}

func (r *Refunds) GetRefundByNo(ctx context.Context, refundNo string) (*entity.RefundRecords, error) {
	var refund entity.RefundRecords
	err := dao.RefundRecords.Ctx(ctx).Where("refund_no", refundNo).Scan(&refund)
	if err != nil {
		return nil, err
	}
	if refund.Id == 0 {
		return nil, gerror.New("退款记录不存在")
	}
	return &refund, nil
}

func generateRefundNo(userId int64) string {
	return fmt.Sprintf("REF%d%d", userId, time.Now().UnixMilli())
}
