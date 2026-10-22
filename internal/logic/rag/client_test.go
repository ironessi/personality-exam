package rag

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"PersonalityTest/internal/consts"
)

func newTestPending(t *testing.T, baseURL string) *Pending {
	t.Helper()
	req, err := BuildRequest(BuildInput{TestType: consts.MBTI, Result: mbtiResult()})
	if err != nil {
		t.Fatalf("BuildRequest 失败: %v", err)
	}
	return &Pending{
		TestType: consts.MBTI,
		Request:  req,
		cfg: Config{
			Enable:              true,
			BaseURL:             baseURL,
			APIKey:              "test-key",
			SyncDeadlineSeconds: 3,
			TestTypes:           []int{consts.MBTI},
		},
	}
}

// cannedResponse 模拟 Python 侧的成功响应（字段与 schemas.py 对齐）。
const cannedResponse = `{
  "success": true,
  "data": {
    "enhanced_text": "你的外倾维度偏高，说明你更倾向于从外部世界获得能量。",
    "sections": [{"title": "维度解读", "content": "外倾者…"}],
    "suggestions": [{"text": "多参与协作式任务", "basis": "外倾者从外部世界获得能量", "chunk_id": "typology_00012", "evidence_quote": "外倾的人思考过程往往需要边说边想", "verified": true}],
    "citations": [{
      "claim": "外倾者从外部世界获得能量",
      "source": "MBTI维度强度解读.md",
      "chunk_id": "typology_00012",
      "evidence_quote": "外倾的人思考过程往往需要边说边想",
      "verified": true,
      "section": 0
    }],
    "model": "deepseek-v4-pro",
    "prompt_tokens": 1674,
    "completion_tokens": 1139,
    "latency_ms": 9250
  },
  "disclaimer": "本结果基于心理学公开文献生成，仅供参考，不构成医学诊断。"
}`

func TestClientEnhanceAndWhitelistStore(t *testing.T) {
	var gotAPIKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("X-API-Key")
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q，期望 application/json", ct)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(cannedResponse))
	}))
	defer srv.Close()

	p := newTestPending(t, srv.URL)
	rag := p.TryInline(context.Background())
	if rag == nil {
		t.Fatal("TryInline 返回 nil，期望拿到增强结果")
	}
	if gotAPIKey != "test-key" {
		t.Errorf("X-API-Key = %q，期望 test-key", gotAPIKey)
	}

	// 业务字段必须在。
	if rag["enhanced_text"] == "" || rag["sections"] == nil || rag["citations"] == nil {
		t.Errorf("落库结构缺业务字段: %v", rag)
	}
	// 历史值是 "ok"（Go 自造），2026-10-21 统一成 Python 枚举里的 "ready"。
	if rag["status"] != "ready" {
		t.Errorf("status = %v，有引用时期望 ok", rag["status"])
	}

	// 成本与模型信息**绝不能**落库（schemas.py 的 EnhanceData 注释明确要求）。
	for _, forbidden := range []string{
		"model", "prompt_tokens", "completion_tokens", "latency_ms",
	} {
		if _, ok := rag[forbidden]; ok {
			t.Errorf("落库结构里出现了 %q，应被白名单挡掉", forbidden)
		}
	}
}

// TestClientKeepsSuggestionEvidence：建议的验证标记必须活着到落库结构里。
//
// 回归用例：`Suggestion` 一度只有 text/basis 两个字段，json.Unmarshal 会把
// chunk_id / evidence_quote / verified 静默丢掉。验证照做、库里没有，
// 等于 ADR-0003 的「建议也走同一套验证」没有落地——而单测与接口都不报错。
func TestClientKeepsSuggestionEvidence(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(cannedResponse))
	}))
	defer srv.Close()

	p := newTestPending(t, srv.URL)
	rag := p.TryInline(context.Background())
	if rag == nil {
		t.Fatal("TryInline 返回 nil")
	}

	suggestions, ok := rag["suggestions"].([]Suggestion)
	if !ok || len(suggestions) != 1 {
		t.Fatalf("suggestions = %#v，期望 1 条 []Suggestion", rag["suggestions"])
	}
	got := suggestions[0]
	if !got.Verified {
		t.Error("verified 丢了：评审人无法区分「有据的建议」与「模型自己编的建议」")
	}
	if got.ChunkID != "typology_00012" {
		t.Errorf("chunk_id = %q，期望 typology_00012", got.ChunkID)
	}
	if got.EvidenceQuote == "" {
		t.Error("evidence_quote 丢了：前端与评审表无法回查依据")
	}
	// 落库后仍要能序列化出这三个键（前端读的是 JSON，不是 Go 结构体）。
	raw, err := json.Marshal(rag)
	if err != nil {
		t.Fatalf("序列化落库结构失败: %v", err)
	}
	for _, key := range []string{`"chunk_id"`, `"evidence_quote"`, `"verified"`} {
		if !strings.Contains(string(raw), key) {
			t.Errorf("落库 JSON 缺 %s: %s", key, raw)
		}
	}
}

// TestClientDegradedWhenNoCitations：一条引用都没过验证时，
// 状态要标成 degraded，不能让它看起来"有依据"。
func TestClientDegradedWhenNoCitations(t *testing.T) {
	body := strings.Replace(cannedResponse, `"citations": [{`, `"citations": [], "_x": [{`, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	p := newTestPending(t, srv.URL)
	rag := p.TryInline(context.Background())
	if rag == nil {
		t.Fatal("TryInline 返回 nil")
	}
	if rag["status"] != "degraded" {
		t.Errorf("status = %v，无引用时期望 degraded", rag["status"])
	}
}

// TestClientReturnsErrorOn422：契约不匹配（如缺 threshold）在这里必须变成
// 可读的错误，而不是一个空结果——否则排查时只能看到"增强没出来"。
func TestClientReturnsErrorOn422(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"detail":[{"loc":["body","measures",0],"msg":"缺少 threshold"}]}`))
	}))
	defer srv.Close()

	p := newTestPending(t, srv.URL)
	if rag := p.TryInline(context.Background()); rag != nil {
		t.Fatalf("422 时 TryInline 应返回 nil，实际: %v", rag)
	}

	// 直接调 client 看错误文本里有没有响应体片段。
	client := NewClient(p.cfg)
	_, err := client.Enhance(context.Background(), p.Request)
	if err == nil {
		t.Fatal("期望返回错误")
	}
	if !strings.Contains(err.Error(), "422") || !strings.Contains(err.Error(), "threshold") {
		t.Errorf("错误信息应包含状态码与响应体片段，实际: %v", err)
	}
}

// TestTryInlineHonoursSyncDeadline：慢服务必须被 3 秒期限切断，
// 否则"提交接口最多等 3 秒"这个承诺就是假的。
func TestTryInlineHonoursSyncDeadline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	}))
	defer srv.Close()

	p := newTestPending(t, srv.URL)
	p.cfg.SyncDeadlineSeconds = 1

	start := time.Now()
	rag := p.TryInline(context.Background())
	elapsed := time.Since(start)

	if rag != nil {
		t.Fatal("超时期望返回 nil")
	}
	if elapsed > 1500*time.Millisecond {
		t.Errorf("耗时 %s，超过 1s 期限太多（期限没生效？）", elapsed)
	}
}

// TestStoredRejectsFailedResponse：success=false 时不应产出任何落库结构。
func TestStoredRejectsFailedResponse(t *testing.T) {
	var resp EnhanceResponse
	if err := json.Unmarshal([]byte(`{"success":false,"error":"增强生成失败"}`), &resp); err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if got := resp.Stored(); got != nil {
		t.Errorf("success=false 时 Stored() 应返回 nil，实际: %v", got)
	}
}
