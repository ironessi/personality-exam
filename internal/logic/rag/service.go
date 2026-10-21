package rag

import (
	"context"

	"PersonalityTest/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
)

// 本文件是**四道闸**：配置开关 → 灰度白名单 → 增强策略 → 测量项定义。
//
// 四道全过才外呼。任何一道不过都返回 nil，让调用方当作"本次没有增强"处理——
// 增强是旁路，**绝不能因为它失败而让测评提交失败**（D1/D2 决策）。

// Pending 是一次"已经决定要增强、但还没发出去"的请求。
//
// 为什么要分成 Prepare 和 TryInline 两步：ADR-0007 要求同步快路径只等 3 秒，
// 超时转异步，而异步回填需要 `record_id`——它在 INSERT 之后才存在。所以
// 调用顺序必须是：Prepare（拼请求）→ 插库 → 拿到 id → TryInline 或 Enqueue。
type Pending struct {
	TestType int
	Request  *EnhanceRequest

	cfg Config
}

// Prepare 走完四道闸并拼出请求。返回 nil 表示本次不做增强（不报错、不打扰用户）。
func Prepare(
	ctx context.Context, testType int, testName string, result map[string]interface{},
) *Pending {
	cfg := LoadConfig(ctx)

	// 闸 1：总开关。默认 false，没显式打开绝不外呼。
	if !cfg.Enable {
		return nil
	}
	// 闸 2：灰度白名单。它是灰度开关，不是能力上限（ADR-0008）。
	if !cfg.Whitelisted(testType) {
		return nil
	}
	// 闸 3：增强策略。disabled 的两个娱乐测评、以及未登记的编码都在这被拦掉。
	if !consts.CanEnhance(testType) {
		return nil
	}
	// 闸 4：测量项定义。没有阈值就没法算 strength，Python 侧会 422。
	if len(consts.MeasuresOf(testType)) == 0 {
		g.Log().Warningf(ctx,
			"RAG: 测评 %d 在灰度白名单里但没有测量项定义（见 internal/consts/rag.go 的 "+
				"TestMeasures 注释，DISC 属于这种情况），本次跳过增强", testType)
		return nil
	}

	req, err := BuildRequest(BuildInput{
		TestType:       testType,
		TestName:       testName,
		Result:         result,
		MaxSuggestions: cfg.MaxSuggestions,
	})
	if err != nil {
		// 拼不出来通常是注册表与评分代码脱节（字段改名）。记 WARNING 让人能查到，
		// 但**不返回错误**：用户不该为服务端的配置漂移买单。
		g.Log().Warningf(ctx, "RAG: 拼增强请求失败，本次跳过: %v", err)
		return nil
	}

	return &Pending{TestType: testType, Request: req, cfg: cfg}
}

// TryInline 在 3 秒硬期限内尝试同步增强（ADR-0007 的快路径）。
//
// 返回 nil 表示"期限内没拿到结果"，调用方应当把请求交给 Enqueue 异步兜底。
// 这里区分不了"超时"和"服务真的挂了"，也不需要区分——两者的处理完全一样。
func (p *Pending) TryInline(ctx context.Context) map[string]interface{} {
	deadline, cancel := context.WithTimeout(ctx, p.cfg.SyncDeadline())
	defer cancel()

	client := NewClient(p.cfg)
	resp, err := client.Enhance(deadline, p.Request)
	if err != nil {
		g.Log().Infof(ctx, "RAG: 测评 %d 同步增强未在期限内完成，转异步: %v", p.TestType, err)
		return nil
	}
	if !resp.Success {
		g.Log().Warningf(ctx, "RAG: 测评 %d 增强失败: %s", p.TestType, resp.Error)
		return nil
	}

	rag := resp.Stored()
	if rag == nil {
		g.Log().Warningf(ctx, "RAG: 测评 %d 增强响应为空", p.TestType)
		return nil
	}
	g.Log().Infof(ctx, "RAG: 测评 %d 同步增强命中，sections=%d citations=%d",
		p.TestType, len(resp.Data.Sections), len(resp.Data.Citations))
	return rag
}
