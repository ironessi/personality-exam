package rag

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
)

// 本文件是「同一份答案不重复花钱」的那一层。
//
// 指纹只覆盖**决定解读内容**的东西：测评类型、各测量项的取值与极性、标志位。
// 刻意不含用户 ID、记录 ID、时间——增强文本描述的是"这份档案"，不是"这个人"，
// 同一档案换个人提交，解读本来就该一样；反过来，把用户 ID 放进指纹会让缓存
// 命中率降到接近 0，缓存也就白做了。
//
// 缓存是**可选加速**，不是正确性依赖：Redis 不可用、反序列化失败、TTL 配 0，
// 一律当作未命中，绝不让提交接口因此失败或变慢。

const cacheKeyPrefix = "rag:cache:"

// cachePayloadVersion 是**结果生成契约**的版本号，只影响指纹，不影响请求契约。
//
// 为什么需要它：指纹原本只覆盖"档案内容"，于是"同一份档案、代码已升级"会命中
// 升级前缓存的旧 payload。实测踩过两次：
//
//  1. 2026-10-22 给 `Suggestion` 补上 chunk_id/evidence_quote/verified 之后，
//     MBTI 两份档案仍然命中旧缓存，库里存的建议依旧只有 {text, basis}，
//     验收脚本报 suggestion_pass_rate=0，看起来像修复没生效。
//  2. 同一天改 `strength_band()` 之后（因子负分不再写成"偏向"），旧缓存继续
//     发"正常范围 + 明显偏向"的自相矛盾文本。
//
// 缓存 TTL 是 30 天，靠过期自愈太久，所以：**凡是会改变生成结果的东西变了就 +1**，
// 包括 `Stored()` 的字段集合、`_SYSTEM_PROMPT`、`strength_band` 一类措辞规则。
// 旧 key 不会被删除，只是再也不会被读到（30 天后自然过期）。
const cachePayloadVersion = 3

// cacheFingerprint 计算请求指纹。
//
// 用结构体而不是 map 来序列化：Go 的 map 键序在 JSON 里是稳定的（encoding/json
// 会对 map 键排序），但字段集合容易漏改；结构体则会在编译期暴露漏字段。
func cacheFingerprint(req *EnhanceRequest) (string, error) {
	return fingerprintWithPayloadVersion(req, cachePayloadVersion)
}

// fingerprintWithPayloadVersion 是 cacheFingerprint 的实现，payload 版本作为参数
// 传入只是为了让测试能证明「版本变了，key 一定跟着变」——否则这条规则只能靠人记。
func fingerprintWithPayloadVersion(req *EnhanceRequest, payloadVersion int) (string, error) {
	if req == nil {
		return "", nil
	}

	type measureFP struct {
		Key         string   `json:"k"`
		Score       float64  `json:"s"`
		Strength    *float64 `json:"st,omitempty"`
		Max         *float64 `json:"mx,omitempty"`
		Threshold   *float64 `json:"th,omitempty"`
		Cutoff      *float64 `json:"co,omitempty"`
		Pole        string   `json:"p,omitempty"`
		FactorLevel string   `json:"fl,omitempty"`
	}

	fp := struct {
		SchemaVersion int         `json:"v"`
		TestCode      int         `json:"t"`
		Measures      []measureFP `json:"m"`
		// needs_referral 由 buildFlags 折进 Flags，不单独占字段。
		Flags []string `json:"f,omitempty"`
		// Payload 是**落库结构**的版本，与请求契约无关：字段集合变了就 +1，
		// 否则旧缓存会把旧形状的结果继续发出去（见 cachePayloadVersion 注释）。
		Payload int `json:"p"`
	}{
		SchemaVersion: req.SchemaVersion,
		TestCode:      req.Test.Code,
		Flags:         req.Flags,
		Payload:       payloadVersion,
	}

	for _, m := range req.Measures {
		fp.Measures = append(fp.Measures, measureFP{
			Key:         m.Key,
			Score:       m.Score,
			Strength:    m.Strength,
			Max:         m.Max,
			Threshold:   m.Threshold,
			Cutoff:      m.ClinicalCutoff,
			Pole:        m.Pole,
			FactorLevel: m.FactorLevel,
		})
	}

	raw, err := json.Marshal(fp)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

// TryCache 读缓存。任何异常都返回 nil（=未命中），不返回错误。
func (p *Pending) TryCache(ctx context.Context) map[string]interface{} {
	if p == nil || p.cfg.CacheTTL() <= 0 {
		return nil
	}
	key, err := cacheFingerprint(p.Request)
	if err != nil || key == "" {
		return nil
	}

	readCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	rdb, err := redisClient(p.cfg)
	if err != nil {
		return nil
	}
	defer rdb.Close(readCtx)

	raw, err := rdb.Get(readCtx, cacheKeyPrefix+key)
	if err != nil || raw.IsNil() || raw.IsEmpty() {
		return nil
	}

	var stored map[string]interface{}
	if err := json.Unmarshal(raw.Bytes(), &stored); err != nil {
		return nil
	}
	return stored
}

// SaveCache 写缓存（失败只丢日志，不影响调用方）。
func (p *Pending) SaveCache(ctx context.Context, stored map[string]interface{}) {
	saveCache(ctx, p.cfg, p.Request, stored)
}

func saveCache(ctx context.Context, cfg Config, req *EnhanceRequest, stored map[string]interface{}) {
	ttl := cfg.CacheTTL()
	if ttl <= 0 || stored == nil {
		return
	}
	key, err := cacheFingerprint(req)
	if err != nil || key == "" {
		// 指纹算不出来就别缓存：宁可多花一次 token，也不能把两份不同档案
		// 的解读写到同一个键上。
		return
	}

	raw, err := json.Marshal(stored)
	if err != nil {
		return
	}

	writeCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	rdb, err := redisClient(cfg)
	if err != nil {
		return
	}
	defer rdb.Close(writeCtx)

	sec := int64(ttl / time.Second)
	if sec < 1 {
		sec = 1
	}
	_, _ = rdb.Set(writeCtx, cacheKeyPrefix+key, string(raw), gredis.SetOption{
		TTLOption: gredis.TTLOption{EX: &sec},
	})
}
