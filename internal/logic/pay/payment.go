package pay

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	notify_request "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/notify/request"
	order_request "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/request"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

var Zero = decimal.NewFromInt(0)

type Payment struct{}

func New() *Payment {
	return &Payment{}
}

// generateOrderNo 生成唯一商户订单号
func generateOrderNo(userId uint64) string {
	return fmt.Sprintf("OT%d%d", userId, time.Now().UnixMilli())
}

// generateTransNo 生成唯一支付流水号
func generateTransNo(userId uint64) string {
	return fmt.Sprintf("TR%d%d", userId, time.Now().UnixMilli())
}

// CreateOrder 创建支付订单并生成支付URL
func (p *Payment) CreateOrder(ctx context.Context, userId, testTypeId int64, payType, terminal, idempotencyKey string) (orderNo, payUrl string, err error) {
	// 1. 幂等令牌校验与锁定
	_, err = dao.IdempotencyTokens.Ctx(ctx).Data(do.IdempotencyTokens{
		Token:        idempotencyKey,
		BusinessType: "create_order",
		BusinessId:   fmt.Sprintf("%d_%d", userId, testTypeId),
		Status:       0,
		ExpireAt:     gtime.Timestamp() + 300, // 5分钟
	}).Insert()
	if err != nil {
		return "", "", gerror.New("订单处理中或存在重复请求，请不要重复点击")
	}

	// 最终根据处理结果更新令牌状态
	defer func() {
		status := 1 // 成功
		if err != nil {
			status = 2 // 失败，释放后可重试 (业务侧可定期清理或根据实际情况提供重试机会)
		}
		dao.IdempotencyTokens.Ctx(ctx).WherePri(idempotencyKey).Data(g.Map{"status": status}).Update()
	}()

	// 2. 查询问卷信息
	var testType entity.TestTypes
	if err = dao.TestTypes.Ctx(ctx).WherePri(testTypeId).Scan(&testType); err != nil {
		return
	}
	if testType.Id == 0 {
		err = gerror.New("问卷不存在")
		return
	}

	// 3. 检查用户是否已支付 (业务订单层面)
	var existOrder entity.TestOrders
	_ = dao.TestOrders.Ctx(ctx).
		Where("user_id", userId).
		Where("test_id", testTypeId).
		Where("status", 1).
		Scan(&existOrder)
	if existOrder.Id > 0 {
		err = gerror.New("您已支付过该问卷报告")
		return
	}

	amount := testType.Price

	var transNo string

	// 3. 开启事务创建订单和流水
	err = dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 3.1 创建业务订单
		orderNo = generateOrderNo(uint64(userId))
		_, err = dao.TestOrders.Ctx(ctx).Data(do.TestOrders{
			OrderNo:       orderNo,
			UserId:        userId,
			TestId:        testTypeId,
			TestName:      testType.Name,
			OriginalPrice: amount,
			PayableAmount: amount,
			Status:        0,
			ExpireAt:      gtime.Timestamp() + 1800, // 30分钟
		}).Insert()
		if err != nil {
			return err
		}

		// 3.2 创建支付流水
		transNo = generateTransNo(uint64(userId))
		_, err = dao.PaymentTransactions.Ctx(ctx).Data(do.PaymentTransactions{
			TransactionNo: transNo,
			OrderNo:       orderNo,
			UserId:        userId,
			Amount:        amount,
			Channel:       payType,
			Status:        0,
			ExpireAt:      gtime.Timestamp() + 1800, // 30分钟
		}).Insert()

		return err
	})
	if err != nil {
		return
	}

	// 4. 根据渠道和终端生成支付链接
	switch payType {
	case "alipay":
		// 获取订单中的结果 ID
		var orderRes entity.TestOrders
		dao.TestOrders.Ctx(ctx).Where("order_no", orderNo).Scan(&orderRes)

		payUrl, err = p.createAlipayOrder(ctx, transNo, testType.Name, decimal.NewFromFloat(amount), terminal, orderRes.ResultReportId)
	case "wechat":
		payUrl, err = p.createWechatOrder(ctx, transNo, testType.Name, decimal.NewFromFloat(amount), terminal)
	default:
		err = gerror.New("暂不支持该支付渠道")
	}

	// 记录操作日志
	p.logAction(ctx, "", orderNo, "CREATE", g.Map{
		"payType": payType, "terminal": terminal, "amount": amount,
	}, g.Map{"payUrl": payUrl}, err)

	return
}

// logAction 记录支付日志
func (p *Payment) logAction(ctx context.Context, transNo, orderNo, action string, req, res interface{}, err error) {
	result := 1
	errMsg := ""
	if err != nil {
		result = 0
		errMsg = err.Error()
	}
	_, _ = dao.PaymentLogs.Ctx(ctx).Data(do.PaymentLogs{
		TransactionNo: transNo,
		OrderNo:       orderNo,
		Action:        action,
		RequestData:   gjson.MustEncodeString(req),
		ResponseData:  gjson.MustEncodeString(res),
		Result:        result,
		ErrorMsg:      errMsg,
	}).Insert()
}

// createAlipayOrder 调用支付宝 OpenAPI 生成支付链接
func (p *Payment) createAlipayOrder(ctx context.Context, outTradeNo, subject string, amount decimal.Decimal, terminal string, resultReportId uint64) (string, error) {
	client, cfg, err := NewAlipayClient(ctx)
	if err != nil {
		return "", gerror.Wrapf(err, "初始化支付宝客户端失败")
	}

	bm := make(gopay.BodyMap)
	returnUrl := cfg.ReturnUrl
	if resultReportId > 0 {
		if !strings.HasSuffix(returnUrl, "/") {
			returnUrl += "/"
		}
		returnUrl += fmt.Sprintf("%d", resultReportId)
	}

	// 清理 subject 中的特殊字符，防止 INVALID_PARAMETER (支付宝要求不可包含 / = & 等)
	subject = strings.NewReplacer("/", "", "=", "", "&", "", "(", "", ")", "", "（", "", "）", "").Replace(subject)

	bm.Set("out_trade_no", outTradeNo).
		Set("total_amount", amount.StringFixed(2)).
		Set("subject", subject+"报告")

	client.SetReturnUrl(returnUrl)

	if terminal == "mobile" {
		bm.Set("product_code", "QUICK_WAP_WAY")
		payUrl, err := client.TradeWapPay(ctx, bm)
		return payUrl, err
	}
	bm.Set("product_code", "FAST_INSTANT_TRADE_PAY").
		Set("qr_pay_mode", "1")

	payUrl, err := client.TradePagePay(ctx, bm)
	return payUrl, err
}

// createWechatOrder 调用微信支付 (PowerWeChat v3)
func (p *Payment) createWechatOrder(ctx context.Context, outTradeNo, subject string, amount decimal.Decimal, terminal string) (string, error) {
	client, _, err := NewWechatClient(ctx)
	if err != nil {
		return "", gerror.Wrapf(err, "初始化微信支付客户端失败")
	}

	totalFee := int(amount.Mul(decimal.NewFromInt(100)).IntPart())

	if terminal == "pc" {
		res, err := client.Order.TransactionNative(ctx, &order_request.RequestNativePrepay{
			Amount: &order_request.NativeAmount{
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
	} else {
		res, err := client.Order.TransactionH5(ctx, &order_request.RequestH5Prepay{
			Amount: &order_request.H5Amount{
				Total:    totalFee,
				Currency: "CNY",
			},
			Description: subject + "分析报告",
			OutTradeNo:  outTradeNo,
			SceneInfo: &order_request.H5SceneInfo{
				PayerClientIP: "127.0.0.1",
				H5Info: &order_request.H5H5Info{
					Type: "Wap",
				},
			},
		})
		if err != nil {
			return "", err
		}
		return res.H5URL, nil
	}
}

// HandleAlipayNotify 处理支付宝异步回调
func (p *Payment) HandleAlipayNotify(ctx context.Context) (err error) {
	r := g.RequestFromCtx(ctx)

	// 获取所有参数用于验签 (GetMap 会自动处理表单解析)
	params := r.GetMap()
	notifyReq := make(gopay.BodyMap)
	for k, v := range params {
		notifyReq.Set(k, v)
	}

	rawBody := r.GetBodyString()
	tradeStatus := notifyReq.GetString("trade_status")
	outTradeNo := notifyReq.GetString("out_trade_no")
	tradeNo := notifyReq.GetString("trade_no")

	g.Log().Debugf(ctx, "支付宝回调原始数据: out_trade_no=%s, trade_status=%s, trade_no=%s", outTradeNo, tradeStatus, tradeNo)

	// 1. 记录原始回调通知
	notifyId, _ := dao.PaymentNotifications.Ctx(ctx).Data(do.PaymentNotifications{
		Channel:       "alipay",
		NotifyType:    "pay",
		RawHeaders:    gjson.MustEncodeString(r.Header),
		RawBody:       rawBody,
		ParsedData:    gjson.MustEncodeString(notifyReq),
		Status:        1, // 处理中
		TransactionNo: tradeNo,
		OrderNo:       outTradeNo,
	}).InsertAndGetId()

	defer func() {
		if notifyId > 0 {
			status := 2 // 成功
			errMsg := ""
			if err != nil {
				status = 3 // 失败
				errMsg = err.Error()
			}
			dao.PaymentNotifications.Ctx(ctx).WherePri(notifyId).Data(do.PaymentNotifications{
				Status:      status,
				ErrorMsg:    errMsg,
				ProcessedAt: gtime.Now(),
			}).Update()
		}
	}()

	_, cfg, err := NewAlipayClient(ctx)
	if err != nil {
		return err
	}

	// 根据配置动态选择验签方式：公钥模式 vs 证书模式
	var ok bool
	if strings.TrimSpace(cfg.AlipayPublicKey) != "" {
		// 公钥模式（沙箱常用）
		g.Log().Debug(ctx, "支付宝回调验签：使用 [公钥模式]")
		ok, err = alipay.VerifySign(cfg.AlipayPublicKey, notifyReq)
	} else {
		// 证书模式（生产常用）
		g.Log().Debug(ctx, "支付宝回调验签：使用 [证书模式]")
		certPath := getAbsPath(ctx, cfg.AlipayPublicCert)
		var alipayPublicCert []byte
		alipayPublicCert, err = os.ReadFile(certPath)
		if err != nil {
			err = gerror.Wrapf(err, "读取支付宝公钥证书失败: %s", certPath)
			return err
		}
		ok, err = alipay.VerifySignWithCert(alipayPublicCert, notifyReq)
	}
	if err != nil || !ok {
		g.Log().Errorf(ctx, "支付宝回调验签失败: ok=%v, err=%v", ok, err)
		err = gerror.New("支付宝回调验签失败")
		return err
	}

	p.logAction(ctx, tradeNo, outTradeNo, "NOTIFY_ALIPAY", notifyReq, nil, nil)

	if tradeStatus != "TRADE_SUCCESS" {
		return nil
	}

	// 开启事务更新流水和订单
	err = dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 更新支付流水 (悲观锁避免并发回调处理冲突)
		var trans entity.PaymentTransactions
		err := dao.PaymentTransactions.Ctx(ctx).
			Where("transaction_no", outTradeNo).
			LockUpdate().
			Scan(&trans)
		if err != nil || trans.Id == 0 {
			return gerror.New("未找到对应支付流水")
		}
		if trans.Status == 2 {
			return nil // 已处理过
		}

		_, err = dao.PaymentTransactions.Ctx(ctx).
			WherePri(trans.Id).
			Data(do.PaymentTransactions{
				Status:        2, // 成功
				AlipayTradeNo: tradeNo,
				PaidAt:        gtime.Now(),
			}).Update()
		if err != nil {
			return err
		}

		// 2. 更新业务订单
		_, err = dao.TestOrders.Ctx(ctx).
			Where("order_no", trans.OrderNo).
			Where("status", 0).
			Data(do.TestOrders{
				Status:            1, // 已支付
				PaidAt:            gtime.Now(),
				PaidChannel:       "alipay",
				PaidTransactionId: tradeNo,
			}).Update()

		return err
	})
	return err
}

// HandleWechatNotify 处理微信支付异步回调
func (p *Payment) HandleWechatNotify(ctx context.Context) error {
	client, _, err := NewWechatClient(ctx)
	if err != nil {
		return err
	}

	r := g.RequestFromCtx(ctx)
	rawBody := r.GetBodyString()

	_, err = client.HandlePaidNotify(r.Request, func(message *notify_request.RequestNotify, transaction *models.Transaction, fail func(message string)) interface{} {
		if transaction != nil {
			// 1. 记录原始回调通知
			notifyId, _ := dao.PaymentNotifications.Ctx(ctx).Data(do.PaymentNotifications{
				Channel:       "wechat",
				NotifyType:    "pay",
				RawHeaders:    gjson.MustEncodeString(r.Header),
				RawBody:       rawBody,
				ParsedData:    gjson.MustEncodeString(transaction),
				Status:        1, // 处理中
				TransactionNo: transaction.TransactionID,
				OrderNo:       transaction.OutTradeNo,
			}).InsertAndGetId()

			processErr := error(nil)
			defer func() {
				if notifyId > 0 {
					status := 2 // 成功
					errMsg := ""
					if processErr != nil {
						status = 3 // 失败
						errMsg = processErr.Error()
					}
					dao.PaymentNotifications.Ctx(ctx).WherePri(notifyId).Data(do.PaymentNotifications{
						Status:      status,
						ErrorMsg:    errMsg,
						ProcessedAt: gtime.Now(),
					}).Update()
				}
			}()

			if transaction.TradeState == "SUCCESS" {
				g.Log().Infof(ctx, "微信支付回调成功: %s", transaction.OutTradeNo)

				processErr = dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
					// 1. 更新支付流水加锁
					var trans entity.PaymentTransactions
					err := dao.PaymentTransactions.Ctx(ctx).
						Where("transaction_no", transaction.OutTradeNo).
						LockUpdate().
						Scan(&trans)
					if err != nil || trans.Id == 0 {
						return gerror.New("未找到对应支付流水")
					}
					if trans.Status == 2 {
						return nil // 已处理过
					}

					_, err = dao.PaymentTransactions.Ctx(ctx).
						WherePri(trans.Id).
						Data(do.PaymentTransactions{
							Status:              2,
							WechatTransactionId: transaction.TransactionID,
							PaidAt:              gtime.Now(),
						}).Update()
					if err != nil {
						return err
					}

					// 2. 更新业务订单
					_, err = dao.TestOrders.Ctx(ctx).
						Where("order_no", trans.OrderNo).
						Where("status", 0).
						Data(do.TestOrders{
							Status:            1,
							PaidAt:            gtime.Now(),
							PaidChannel:       "wechat",
							PaidTransactionId: transaction.TransactionID,
						}).Update()
					return err
				})
				if processErr != nil {
					fail(processErr.Error())
				}
			}
		}
		return nil
	})

	return err
}

// CheckPaid 检查用户是否已为某问卷支付
func (p *Payment) CheckPaid(ctx context.Context, userId, testTypeId int64) (bool, error) {
	count, err := dao.TestOrders.Ctx(ctx).
		Where("user_id", userId).
		Where("test_id", testTypeId).
		Where("status", 1).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// QueryOrder 查询订单状态
func (p *Payment) QueryOrder(ctx context.Context, orderNo string) (*entity.TestOrders, error) {
	var order entity.TestOrders
	err := dao.TestOrders.Ctx(ctx).Where("order_no", orderNo).Scan(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// RefundOrder 申请退款并全额退回
func (p *Payment) RefundOrder(ctx context.Context, userId int64, orderNo, reason string) (string, error) {
	var order entity.TestOrders
	err := dao.TestOrders.Ctx(ctx).
		Where("order_no", orderNo).
		Where("user_id", userId).
		Scan(&order)
	if err != nil || order.Id == 0 {
		return "", gerror.New("订单不存在或不属于该用户")
	}
	if order.Status != 1 && order.Status != 3 {
		return "", gerror.New("当前订单状态不支持退款")
	}

	var trans entity.PaymentTransactions
	err = dao.PaymentTransactions.Ctx(ctx).
		Where("order_no", orderNo).
		Where("status", 2). // 支付成功
		Scan(&trans)
	if err != nil || trans.Id == 0 {
		return "", gerror.New("未找到成功的支付流水")
	}

	refundNo := fmt.Sprintf("RF%d%d", userId, time.Now().UnixMilli())

	// 记录退款流水
	refundId, err := dao.RefundRecords.Ctx(ctx).Data(do.RefundRecords{
		RefundNo:      refundNo,
		TransactionNo: trans.TransactionNo,
		OrderNo:       orderNo,
		RefundAmount:  trans.Amount,
		Reason:        reason,
		Status:        1, // 处理中
	}).InsertAndGetId()
	if err != nil {
		return "", err
	}

	defer func() {
		// 根据结果更新状态
		status := 2 // 成功
		if err != nil {
			status = 3 // 失败
		}
		dao.RefundRecords.Ctx(ctx).WherePri(refundId).Data(g.Map{"status": status}).Update()

		if err == nil {
			// 更新订单和流水状态
			dao.TestOrders.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				dao.TestOrders.Ctx(ctx).Where("order_no", orderNo).Data(g.Map{"status": 4}).Update()
				dao.PaymentTransactions.Ctx(ctx).Where("order_no", orderNo).Where("status", 2).Data(g.Map{"status": 5}).Update()
				return nil
			})
		}
	}()

	if trans.Channel == "alipay" {
		client, _, errClient := NewAlipayClient(ctx)
		if errClient != nil {
			err = errClient
			return "", err
		}
		bm := make(gopay.BodyMap)
		bm.Set("out_trade_no", trans.TransactionNo).
			Set("refund_amount", fmt.Sprintf("%.2f", trans.Amount)).
			Set("refund_reason", reason).
			Set("out_request_no", refundNo)

		res, errRefund := client.TradeRefund(ctx, bm)
		if errRefund != nil {
			err = errRefund
			return "", err
		}
		if res.Response.Code != "10000" {
			err = gerror.New(res.Response.Msg + " " + res.Response.SubMsg)
			return "", err
		}

	} else if trans.Channel == "wechat" {
		// TODO: PowerWeChat v3 退款 API 集成 (需要商户证书)
		// 暂作占位，或可直接抛错。
		err = gerror.New("微信退款功能正在开发中，暂不开放")
		return "", err
	} else {
		err = gerror.New("不支持的支付渠道")
		return "", err
	}

	return refundNo, nil
}

// DailyReconciliation 每日对账调度任务
func (p *Payment) DailyReconciliation(ctx context.Context) {
	g.Log().Info(ctx, "开始执行每日支付对账任务...")

	// 1. 创建当日对账单记录 reconciliation_records
	// 2. 调用支付宝/微信API获取昨日账单
	// 3. 将解析的数据入库至 reconciliation_channel_data
	// 4. 与本地 payment_transactions 进行比对
	// 5. 发现差异记录入 reconciliation_exceptions 表以供人工处理

	g.Log().Info(ctx, "对账任务执行完毕 (Skeleton)")
}
