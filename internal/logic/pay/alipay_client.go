package pay

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type AlipayConfig struct {
	AppId              string `yaml:"app_id"`
	PrivateKey         string `yaml:"private_key"`
	AppCertPath        string `yaml:"app_cert_path"`
	AlipayRootCertPath string `yaml:"alipay_root_cert_path"`
	AlipayPublicCert   string `yaml:"alipay_public_cert_path"`
	AlipayPublicKey    string `yaml:"alipay_public_key"` // 公钥模式专用
	EncryptKey         string `yaml:"encrypt_key"`      // AES加密密钥
	NotifyUrl          string `yaml:"notify_url"`
	ReturnUrl          string `yaml:"return_url"`
	IsProd             bool   `yaml:"is_prod"`
}

func getProjectRoot() string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(execPath)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func getAbsPath(ctx context.Context, relativePath string) string {
	if relativePath == "" {
		return ""
	}
	relativePath = filepath.ToSlash(relativePath)
	if filepath.IsAbs(relativePath) {
		return relativePath
	}

	// 1. Check current working directory
	workDir, err := os.Getwd()
	if err == nil {
		workDir = filepath.ToSlash(workDir)
		absPath := filepath.Join(workDir, relativePath)
		absPath = filepath.ToSlash(absPath)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// 2. Check project root (where go.mod is)
	root := getProjectRoot()
	if root != "" {
		root = filepath.ToSlash(root)
		absPath := filepath.Join(root, relativePath)
		absPath = filepath.ToSlash(absPath)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// 3. Fallback: check executable directory
	// Useful for server deployments where binary and certs/ are together but no go.mod
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		execDir = filepath.ToSlash(execDir)
		absPath := filepath.Join(execDir, relativePath)
		absPath = filepath.ToSlash(absPath)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return filepath.ToSlash(relativePath)
}

func NewAlipayClient(ctx context.Context) (*alipay.Client, *AlipayConfig, error) {
	var cfg AlipayConfig
	if err := g.Cfg("payment").MustGet(ctx, "alipay").Scan(&cfg); err != nil {
		return nil, nil, err
	}

	client, err := alipay.NewClient(cfg.AppId, cfg.PrivateKey, cfg.IsProd)
	if err != nil {
		return nil, nil, err
	}

	// 如果是非生产环境（沙箱），开启 Debug 模式
	if !cfg.IsProd {
		client.DebugSwitch = gopay.DebugOn
	}

	// 如果提供了 AES 密钥，则设置
	if cfg.EncryptKey != "" {
		client.SetAESKey(cfg.EncryptKey)
	}

	// --- 核心逻辑：区分“公钥模式”还是“证书模式” ---
	if strings.TrimSpace(cfg.AlipayPublicKey) != "" {
		// 1. 公钥模式 (沙箱常用)
		g.Log().Debug(ctx, "支付宝初始化：进入 [公钥模式]")
		client.AutoVerifySign([]byte(cfg.AlipayPublicKey))
	} else {
		// 2. 证书模式 (生产常用)
		g.Log().Debug(ctx, "支付宝初始化：进入 [证书模式]")

		// 补全路径字段
		if strings.TrimSpace(cfg.AlipayPublicCert) == "" {
			cfg.AlipayPublicCert = g.Cfg("payment").MustGet(ctx, "alipay.alipay_public_cert_path").String()
		}
		if strings.TrimSpace(cfg.AppCertPath) == "" {
			cfg.AppCertPath = g.Cfg("payment").MustGet(ctx, "alipay.app_cert_path").String()
		}
		if strings.TrimSpace(cfg.AlipayRootCertPath) == "" {
			cfg.AlipayRootCertPath = g.Cfg("payment").MustGet(ctx, "alipay.alipay_root_cert_path").String()
		}

		appCertPath := getAbsPath(ctx, cfg.AppCertPath)
		rootCertPath := getAbsPath(ctx, cfg.AlipayRootCertPath)
		alipayPublicCert := getAbsPath(ctx, cfg.AlipayPublicCert)

		if strings.Contains(appCertPath, "从支付宝") || strings.Contains(rootCertPath, "从支付宝") {
			return nil, nil, gerror.New("支付宝证书路径配置错误")
		}

		for _, path := range []string{appCertPath, rootCertPath, alipayPublicCert} {
			if path == "" {
				return nil, nil, gerror.New("支付宝证书路径为空")
			}
			if _, err := os.Stat(path); err != nil {
				g.Log().Error(ctx, "支付宝证书文件不存在:", path)
				return nil, nil, gerror.New("支付宝证书文件不存在: " + path)
			}
		}

		if err = client.SetCertSnByPath(appCertPath, rootCertPath, alipayPublicCert); err != nil {
			return nil, nil, err
		}
	}

	client.SetNotifyUrl(cfg.NotifyUrl)

	return client, &cfg, nil
}

