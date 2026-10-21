// Package rag 是 Go 侧对 Python 增强服务的唯一调用入口。
//
// 这一层只做三件事：**拼规范请求**（request.go）、**发 HTTP**（client.go）、
// **按策略决定调不调**（service.go）。它不解析测评结果、不判断构念域——
// 域与策略来自 internal/consts 的注册表（ADR-0008），Python 侧不判断测评类型。
package rag

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Config 是 manifest/config/config.yaml 里 rag: 段的只读投影。
//
// 为什么不定义成配置结构体再绑定：本项目的配置读取一律走
// g.Cfg().MustGet(ctx, "a.b")（见 internal/logic/sms/sms.go:92），
// 没有统一的配置 struct。多引入一套绑定机制只会让"改哪个文件生效"更难查。
type Config struct {
	Enable              bool
	BaseURL             string
	APIKey              string
	SyncDeadlineSeconds int
	TimeoutSeconds      int
	MaxSuggestions      int
	TestTypes           []int

	// 异步兜底用的 Redis（ADR-0007）。本机独立 db=5，与业务缓存隔离。
	RedisAddress string
	RedisDB      int
	RedisPass    string
	QueueKey     string
}

// 默认值刻意与 config.yaml 的注释一致：即使配置项缺失也不会 panic，
// 但 Enable 默认 false——**没有显式打开就绝不外呼**。
const (
	defaultBaseURL             = "http://127.0.0.1:8001"
	defaultSyncDeadlineSeconds = 3
	defaultTimeoutSeconds      = 30
	defaultMaxSuggestions      = 5
	defaultRedisAddress        = "127.0.0.1:6379"
	defaultRedisDB             = 5
	defaultQueueKey            = "rag:queue"
)

// LoadConfig 从配置中心读 rag 段。缺失的键取保守默认值。
func LoadConfig(ctx context.Context) Config {
	cfg := Config{
		Enable:              g.Cfg().MustGet(ctx, "rag.enable", false).Bool(),
		BaseURL:             g.Cfg().MustGet(ctx, "rag.baseUrl", defaultBaseURL).String(),
		APIKey:              g.Cfg().MustGet(ctx, "rag.apiKey", "").String(),
		SyncDeadlineSeconds: g.Cfg().MustGet(ctx, "rag.syncDeadlineSeconds", defaultSyncDeadlineSeconds).Int(),
		TimeoutSeconds:      g.Cfg().MustGet(ctx, "rag.timeoutSeconds", defaultTimeoutSeconds).Int(),
		MaxSuggestions:      g.Cfg().MustGet(ctx, "rag.maxSuggestions", defaultMaxSuggestions).Int(),
	}
	cfg.TestTypes = g.Cfg().MustGet(ctx, "rag.testTypes").Ints()
	cfg.RedisAddress = g.Cfg().MustGet(ctx, "rag.redis.address", defaultRedisAddress).String()
	cfg.RedisDB = g.Cfg().MustGet(ctx, "rag.redis.db", defaultRedisDB).Int()
	cfg.RedisPass = g.Cfg().MustGet(ctx, "rag.redis.pass", "").String()
	cfg.QueueKey = g.Cfg().MustGet(ctx, "rag.queueKey", defaultQueueKey).String()
	return cfg
}

// SyncDeadline 是提交接口内同步尝试增强的硬期限（ADR-0007）。
//
// 实测单次增强（检索 + DeepSeek + 两级验证）在开发机上是 4.8~9.3 秒，
// **远超 3 秒**。所以这里必须是一个真的 deadline，而不是"HTTP 超时"：
// 到点就放弃增强、照常返回测评结果，绝不让用户为增强等待。
func (c Config) SyncDeadline() time.Duration {
	if c.SyncDeadlineSeconds <= 0 {
		return time.Duration(defaultSyncDeadlineSeconds) * time.Second
	}
	return time.Duration(c.SyncDeadlineSeconds) * time.Second
}

// WorkerTimeout 是后台 worker 调 Python 的超时（本切片未实现 worker）。
func (c Config) WorkerTimeout() time.Duration {
	if c.TimeoutSeconds <= 0 {
		return time.Duration(defaultTimeoutSeconds) * time.Second
	}
	return time.Duration(c.TimeoutSeconds) * time.Second
}

// Whitelisted 报告该测评是否在灰度白名单里。
//
// 白名单是**灰度开关**，不是能力上限：42 个测评都有构念域与策略，
// 扩大它只需要改配置，不需要改检索代码（ADR-0008）。
func (c Config) Whitelisted(testType int) bool {
	for _, v := range c.TestTypes {
		if v == testType {
			return true
		}
	}
	return false
}
