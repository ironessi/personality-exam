package pay

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"
	"fmt"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/request"
	"github.com/go-pay/gopay"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

type Qrcode struct{}

func NewQrcode() *Qrcode {
	return &Qrcode{}
}

type CreateQrcodeInput struct {
	OrderNo string
	Channel string
}

type CreateQrcodeOutput struct {
	QrCodeId      string      `json:"qrCodeId"`
	TransactionNo string      `json:"transactionNo"`
	QrType        int         `json:"qrType"`
	QrContent     string      `json:"qrContent"`
	QrImageUrl    string      `json:"qrImageUrl"`
	Amount        float64     `json:"amount"`
	ExpireAt      int64       `json:"expireAt"`
	ValidSeconds  int         `json:"validSeconds"`
}

func (q *Qrcode) CreateQrcode(ctx context.Context, in CreateQrcodeInput) (*CreateQrcodeOutput, error) {
	var order entity.TestOrders
	err := dao.TestOrders.Ctx(ctx).Where("order_no", in.OrderNo).Scan(&order)
	if err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, gerror.New("订单不存在")
	}
	if order.Status != 0 {
		return nil, gerror.New("订单已支付")
	}

	// 强效修复：直接从数据库读取原始值 (Unix 时间戳)
	expireAtVal, _ := dao.TestOrders.Ctx(ctx).Where("order_no", in.OrderNo).Value("expire_at")
	if expireAtVal.IsEmpty() {
		return nil, gerror.New("订单过期时间读取失败")
	}

	expireAtTs := expireAtVal.Int64()
	if expireAtTs < gtime.Timestamp() {
		g.Log().Warningf(ctx, "判定过期详情: OrderNo=%s, DB原始值=%d, 系统当前戳=%d",
			in.OrderNo, expireAtTs, gtime.Timestamp())
		return nil, gerror.New("订单已过期")
	}

	var validTrans entity.PaymentTransactions
	var existQr entity.PaymentQrcodes

	// 为了复用未过期的二维码，我们需要先找到对应的 transaction
	var pendingTrans []entity.PaymentTransactions
	dao.PaymentTransactions.Ctx(ctx).Where("order_no", in.OrderNo).Where("status", 0).Scan(&pendingTrans)

	for _, t := range pendingTrans {
		var qr entity.PaymentQrcodes
		dao.PaymentQrcodes.Ctx(ctx).Where("transaction_no", t.TransactionNo).Order("id desc").Limit(1).Scan(&qr)
		if qr.Id > 0 && qr.Status == 0 && qr.ExpireAt > gtime.Timestamp() && t.Channel == in.Channel {
			existQr = qr
			validTrans = t
			break
		}
	}

	transNo := validTrans.TransactionNo
	if validTrans.Id == 0 {
		transNo = generateTransNo(order.UserId)
		transExpireAt := gtime.Timestamp() + 300

		err = dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 直接创建全新的流水
			_, err = dao.PaymentTransactions.Ctx(ctx).Data(do.PaymentTransactions{
				TransactionNo: transNo,
				OrderNo:       in.OrderNo,
				UserId:        order.UserId,
				Amount:        order.PayableAmount,
				Channel:       in.Channel,
				Status:        0,
				ExpireAt:      transExpireAt,
			}).Insert()
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	if existQr.Id > 0 {
		return &CreateQrcodeOutput{
			QrCodeId:      fmt.Sprintf("QR%d", existQr.Id),
			TransactionNo: transNo,
			QrType:        existQr.QrType,
			QrContent:     existQr.QrContent,
			QrImageUrl:    existQr.QrImageUrl,
			Amount:        order.PayableAmount,
			ExpireAt:      existQr.ExpireAt,
			ValidSeconds:  int(existQr.ValidSeconds),
		}, nil
	}

	var qrType int
	var qrContent string
	validSeconds := 300

	switch in.Channel {
	case "alipay":
		qrType = 1
		qrContent, err = q.createAlipayQrcode(ctx, transNo, order.TestName, order.PayableAmount, transNo, order.ResultReportId)
	case "wechat":
		qrType = 2
		qrContent, err = q.createWechatQrcode(ctx, transNo, order.TestName, order.PayableAmount, transNo)
	default:
		return nil, gerror.New("支付渠道不支持")
	}
	if err != nil {
		return nil, err
	}

	qrExpireAt := gtime.Timestamp() + 300

	var qrId int64
	if existQr.Id > 0 {
		// 更新已有记录（例如已过期/失效）而不是再次插入，避免唯一索引冲突
		_, err = dao.PaymentQrcodes.Ctx(ctx).WherePri(existQr.Id).Data(do.PaymentQrcodes{
			QrType:       qrType,
			QrContent:    qrContent,
			Status:       0,
			ExpireAt:     qrExpireAt,
			ValidSeconds: uint(validSeconds),
		}).Update()
		if err != nil {
			return nil, err
		}
		qrId = int64(existQr.Id)
	} else {
		var insertedId int64
		insertedId, err = dao.PaymentQrcodes.Ctx(ctx).Data(do.PaymentQrcodes{
			TransactionNo: transNo,
			QrType:        qrType,
			QrContent:     qrContent,
			Status:        0,
			ExpireAt:      qrExpireAt,
			ValidSeconds:  uint(validSeconds),
		}).InsertAndGetId()
		if err != nil {
			return nil, err
		}
		qrId = insertedId
	}

	return &CreateQrcodeOutput{
		QrCodeId:      fmt.Sprintf("QR%d", qrId),
		TransactionNo: transNo,
		QrType:        qrType,
		QrContent:     qrContent,
		Amount:        order.PayableAmount,
		ExpireAt:      qrExpireAt,
		ValidSeconds:  validSeconds,
	}, nil
}

func (q *Qrcode) createAlipayQrcode(ctx context.Context, outTradeNo, subject string, amount float64, transNo string, resultReportId uint64) (string, error) {
	client, cfg, err := NewAlipayClient(ctx)
	if err != nil {
		return "", gerror.Wrapf(err, "初始化支付宝客户端失败")
	}

	// 清理 subject 中的特殊字符，防止 INVALID_PARAMETER (支付宝要求不可包含 / = & 等)
	subject = strings.NewReplacer("/", "", "=", "", "&", "", "(", "", ")", "", "（", "", "）", "").Replace(subject)

	bm := make(gopay.BodyMap)
	// 三方 out_trade_no 统一使用流水号，保证每次都能生成独立的支付链接，避免 TOTAL_FEE_EXCEED 等冲突报错。
	bm.Set("out_trade_no", outTradeNo).
		Set("total_amount", decimal.NewFromFloat(amount).StringFixed(2)).
		Set("subject", subject+"报告").
		Set("product_code", "FAST_INSTANT_TRADE_PAY").
		Set("qr_pay_mode", "1").
		Set("timeout_express", "5m")

	// 设置同步跳转地址，拼接 resultReportId
	returnUrl := cfg.ReturnUrl
	if returnUrl != "" {
		if resultReportId > 0 {
			if !strings.HasSuffix(returnUrl, "/") {
				returnUrl += "/"
			}
			returnUrl += fmt.Sprintf("%d", resultReportId)
		}
		client.SetReturnUrl(returnUrl)
	}

	// 使用 TradePagePay (电脑网站支付) — 返回支付宝收银台页面URL，页面自带扫码功能
	payUrl, err := client.TradePagePay(ctx, bm)
	if err != nil {
		return "", gerror.Wrapf(err, "支付宝电脑网站支付下单失败")
	}

	return payUrl, nil
}


func (q *Qrcode) createWechatQrcode(ctx context.Context, outTradeNo, subject string, amount float64, transNo string) (string, error) {
	client, _, err := NewWechatClient(ctx)
	if err != nil {
		return "", gerror.Wrapf(err, "初始化微信支付客户端失败")
	}

	totalFee := int(decimal.NewFromFloat(amount).Mul(decimal.NewFromInt(100)).IntPart())

	res, err := client.Order.TransactionNative(ctx, &request.RequestNativePrepay{
		Amount: &request.NativeAmount{
			Total:    totalFee,
			Currency: "CNY",
		},
		Description: subject + "分析报告",
		OutTradeNo:  outTradeNo,
	})
	if err != nil {
		return "", err
	}
	return res.CodeURL, nil
}

type GetQrcodeStatusInput struct {
	QrCodeId      string
	TransactionNo string
}

type GetQrcodeStatusOutput struct {
	QrStatus       int         `json:"qrStatus"`
	PayStatus      int         `json:"payStatus"`
	ScannedAt      *gtime.Time `json:"scannedAt"`
	PaidAt         *gtime.Time `json:"paidAt"`
	ExpireAt       int64       `json:"expireAt"`
	FrontStatus    string      `json:"frontStatus"`
	ResultReportId int64       `json:"resultReportId"`
}

func (q *Qrcode) GetQrcodeStatus(ctx context.Context, in GetQrcodeStatusInput) (*GetQrcodeStatusOutput, error) {
	var qr entity.PaymentQrcodes
	err := dao.PaymentQrcodes.Ctx(ctx).Where("transaction_no", in.TransactionNo).Scan(&qr)
	if err != nil {
		return nil, err
	}
	if qr.Id == 0 {
		return nil, gerror.New("二维码不存在")
	}
	if qr.ExpireAt == 0 {
		return nil, gerror.New("二维码过期时间缺失")
	}

	var trans entity.PaymentTransactions
	dao.PaymentTransactions.Ctx(ctx).Where("transaction_no", in.TransactionNo).Scan(&trans)

	frontStatus := "PENDING"
	if qr.ExpireAt < gtime.Timestamp() {
		frontStatus = "EXPIRED"
	} else if trans.Status == 2 {
		frontStatus = "SUCCESS"
	} else if qr.Status == 1 {
		frontStatus = "SCANNED"
	}

	// 主动对账：如果本地状态仍为 PENDING，则主动向支付宝查询真实支付状态
	if frontStatus == "PENDING" && trans.Id > 0 && trans.Channel == "alipay" {
		synced := q.syncAlipayTradeStatus(ctx, trans.TransactionNo, trans.OrderNo)
		if synced {
			// 重新读取更新后的流水状态
			dao.PaymentTransactions.Ctx(ctx).Where("transaction_no", in.TransactionNo).Scan(&trans)
			if trans.Status == 2 {
				frontStatus = "SUCCESS"
			}
		}
	}

	// 获取关联的结果报告 ID
	var resultReportId int64
	if trans.Id > 0 {
		var order entity.TestOrders
		dao.TestOrders.Ctx(ctx).Where("order_no", trans.OrderNo).Scan(&order)
		resultReportId = int64(order.ResultReportId)
	}

	return &GetQrcodeStatusOutput{
		QrStatus:       qr.Status,
		PayStatus:      trans.Status,
		ScannedAt:      qr.ScannedAt,
		PaidAt:         trans.PaidAt,
		ExpireAt:       qr.ExpireAt,
		FrontStatus:    frontStatus,
		ResultReportId: resultReportId,
	}, nil
}

// syncAlipayTradeStatus 主动向支付宝查询订单状态，如果已支付则同步更新本地数据库
func (q *Qrcode) syncAlipayTradeStatus(ctx context.Context, transNo, orderNo string) bool {
	client, _, err := NewAlipayClient(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "主动对账：初始化支付宝客户端失败: %v", err)
		return false
	}

	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", transNo)

	aliRsp, err := client.TradeQuery(ctx, bm)
	if err != nil {
		g.Log().Warningf(ctx, "主动对账：查询支付宝订单失败: %v", err)
		return false
	}

	if aliRsp.Response.Code != "10000" {
		g.Log().Infof(ctx, "主动对账：支付宝返回非成功: code=%s, msg=%s", aliRsp.Response.Code, aliRsp.Response.SubMsg)
		return false
	}

	tradeStatus := aliRsp.Response.TradeStatus
	tradeNo := aliRsp.Response.TradeNo
	g.Log().Infof(ctx, "主动对账：支付宝查询结果 orderNo=%s, tradeStatus=%s, tradeNo=%s", orderNo, tradeStatus, tradeNo)

	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		return false
	}

	// 支付宝已确认支付成功，更新本地数据库
	err = dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新支付流水
		_, err := dao.PaymentTransactions.Ctx(ctx).
			Where("transaction_no", transNo).
			Where("status", 0).
			Data(do.PaymentTransactions{
				Status:        2, // 支付成功
				AlipayTradeNo: tradeNo,
				PaidAt:        gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}

		// 更新业务订单
		_, err = dao.TestOrders.Ctx(ctx).
			Where("order_no", orderNo).
			Where("status", 0).
			Data(do.TestOrders{
				Status:            1, // 已支付
				PaidAt:            gtime.Now(),
				PaidChannel:       "alipay",
				PaidTransactionId: tradeNo,
			}).Update()
		return err
	})

	if err != nil {
		g.Log().Errorf(ctx, "主动对账：更新本地订单状态失败: %v", err)
		return false
	}

	g.Log().Infof(ctx, "主动对账：成功同步订单状态 orderNo=%s", orderNo)
	return true
}
