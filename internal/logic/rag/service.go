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

// TryInline 尝试在提交线程里就拿到增强结果（ADR-0007 的快路径）。
//
// 返回 nil 表示"这一次没有现成的结果"，调用方应当把请求交给 Enqueue 异步兜底。
//
// 两条出口，顺序不能换：
//
//  1. **缓存命中**：同一份答案档案重复提交（用户刷新、重测同题、灰度验收）不该
//     再花三五秒和一次 token。实测命中时提交接口 39ms 就能带 rag 返回。
//
//  2. **缓存未命中**：只有 `rag.syncOnMiss = true` 才会真的外呼一次。
//     默认 **false**，因为 2026-10-22 的 20 条验收给出了反证：`sync_hit = 0`，
//     而且日志里每个用例都有**两条** verify 记录——第一条来自这里的同步尝试
//     （3 秒期限到、Go 放弃，但 Python 侧的生成不会因此停下），第二条来自 worker
//     重跑。结果是每次提交付两倍 token，用户还在提交页白等 3 秒。
//
//     DeepSeek 生成一份解读要 4~10 秒，3 秒期限几乎必然到期；把期限调大到能覆盖
//     生成时间，就等于要求用户在提交页等十秒。**入队才是正常路径，同步只是缓存
//     命中时的赠品**。
func (p *Pending) TryInline(ctx context.Context) map[string]interface{} {
	// 出口 1：缓存。
	if cached := p.TryCache(ctx); cached != nil {
		g.Log().Infof(ctx, "RAG: 测评 %d 命中缓存，跳过外呼", p.TestType)
		return cached
	}

	// 出口 2：未命中缓存且配置为直接入队 → 立刻返回，让调用方 Enqueue。
	if !p.cfg.SyncOnMiss {
		g.Log().Infof(ctx, "RAG: 测评 %d 未命中缓存，直接入队（syncOnMiss=false）", p.TestType)
		return nil
	}

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

	p.SaveCache(ctx, rag)
	return rag
}
