package pay

import (
	"context"
	"os"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment"
	"github.com/gogf/gf/v2/frame/g"
)

// WechatConfig holds configuration for WeChat Pay client (PowerWeChat).
type WechatConfig struct {
	AppId       string `yaml:"app_id"`
	MchId       string `yaml:"mch_id"`
	MchApiV3Key string `yaml:"mch_v3_key"`
	KeyPath     string `yaml:"key_path"`
	CertPath    string `yaml:"cert_path"`
	SerialNo    string `yaml:"serial_no"` // 商户证书序列号
	NotifyUrl   string `yaml:"notify_url"`
}

// NewWechatClient 初始化微信支付客户端 (PowerWeChat)
func NewWechatClient(ctx context.Context) (*payment.Payment, *WechatConfig, error) {
	var cfg WechatConfig
	if err := g.Cfg("payment").MustGet(ctx, "wechat").Scan(&cfg); err != nil {
		return nil, nil, err
	}

	// 使用统一的绝对路径解析
	absKeyPath := getAbsPath(ctx, cfg.KeyPath)
	absCertPath := getAbsPath(ctx, cfg.CertPath)

	// 读取 API 证书私钥
	keyContent, err := os.ReadFile(absKeyPath)
	if err != nil {
		return nil, nil, err
	}

	app, err := payment.NewPayment(&payment.UserConfig{
		AppID:       cfg.AppId,
		MchID:       cfg.MchId,
		MchApiV3Key: cfg.MchApiV3Key,
		Key:         string(keyContent),
		CertPath:    absCertPath,
		KeyPath:     absKeyPath,
		SerialNo:    cfg.SerialNo,
		NotifyURL:   cfg.NotifyUrl,
	})

	return app, &cfg, err
}
