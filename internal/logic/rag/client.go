package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// 本文件只负责"把请求发出去、把响应读回来"，不判断该不该发（那是 service.go）。
//
// 用标准库 net/http 而不是 ghttp.Client：ghttp 的客户端 API 在 gf 2.x
// 各小版本之间动过签名，而这一层只需要 POST 一个 JSON，标准库更稳、更好测
// （测试里可以直接塞一个 httptest.Server）。

// maxResponseBytes 是响应体读取上限。1 MiB 足够放下最长的增强结果
// （实测 MBTI 约 7 KB），超出说明对面出问题了，不该把内存吃穿。
const maxResponseBytes = 1 << 20

// Client 是 Python 增强服务的 HTTP 客户端。
type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

// NewClient 用配置构造客户端。超时由调用方的 context 决定，
// 这里不设 http.Client.Timeout——两处都设会让"到底是哪个超时"变得不可查。
func NewClient(cfg Config) *Client {
	return &Client{
		BaseURL: strings.TrimRight(cfg.BaseURL, "/"),
		APIKey:  cfg.APIKey,
		HTTP:    &http.Client{},
	}
}

// EnhanceResponse 是 POST /api/v1/enhance 的响应体。
//
// 字段与 rag-service/app/api/schemas.py 的 EnhanceResponse 一一对应。
type EnhanceResponse struct {
	Success    bool         `json:"success"`
	Data       *EnhanceData `json:"data"`
	Error      string       `json:"error"`
	Disclaimer string       `json:"disclaimer"`
}

// EnhanceData 是增强结果本体。
//
// ⚠️ model / prompt_tokens / completion_tokens / latency_ms **不得落库**，
// 也不要回给前端（schemas.py 的注释明确要求 Go 侧显式构造白名单结构体）。
// 所以它们只在这一层存在，Stored() 不会带出去。
type EnhanceData struct {
	EnhancedText string       `json:"enhanced_text"`
	Sections     []Section    `json:"sections"`
	Suggestions  []Suggestion `json:"suggestions"`
	Citations    []Citation   `json:"citations"`

	Model            string `json:"model"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	LatencyMs        int    `json:"latency_ms"`
}

// Section 是一段带小标题的解读正文。
type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Suggestion 是一条建议，basis 指向它的依据块。
type Suggestion struct {
	Text  string `json:"text"`
	Basis string `json:"basis"`
}

// Citation 是一条已通过两级验证的引用。
type Citation struct {
	Claim         string  `json:"claim"`
	Source        string  `json:"source"`
	ChunkID       string  `json:"chunk_id"`
	Score         float64 `json:"score"`
	EvidenceQuote string  `json:"evidence_quote"`
	Verified      bool    `json:"verified"`
	Section       *int    `json:"section"`
}

// Enhance 调一次增强。ctx 的 deadline 就是调用方给的期限，
// 这里不额外包一层——否则 3 秒硬期限会被 30 秒 worker 超时悄悄覆盖掉。
func (c *Client) Enhance(ctx context.Context, req *EnhanceRequest) (*EnhanceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化增强请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, c.BaseURL+"/api/v1/enhance", bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("构造增强请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		httpReq.Header.Set("X-API-Key", c.APIKey)
	}

	start := time.Now()
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用增强服务失败（耗时 %s）: %w", time.Since(start).Round(time.Millisecond), err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("读取增强响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// 422 尤其常见：契约不匹配（比如 measures 少了 threshold）。
		// 把响应体前 512 字节带出来，否则只能看到"422"三个字。
		snippet := string(raw)
		if len(snippet) > 512 {
			snippet = snippet[:512] + "…"
		}
		return nil, fmt.Errorf("增强服务返回 %d: %s", resp.StatusCode, snippet)
	}

	var out EnhanceResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("解析增强响应失败: %w（原始响应前 200 字节: %.200s）", err, raw)
	}
	return &out, nil
}

// Stored 把响应转成**允许落库**的白名单结构。
//
// 只保留业务字段 + 免责声明。成本与模型信息留在服务端日志里，
// 落库会让 result_data 变胖并且泄露内部实现（§5.2）。
func (r *EnhanceResponse) Stored() map[string]interface{} {
	if r == nil || !r.Success || r.Data == nil {
		return nil
	}
	d := r.Data
	rag := map[string]interface{}{
		"enhanced_text": d.EnhancedText,
		"sections":      d.Sections,
		"suggestions":   d.Suggestions,
		"citations":     d.Citations,
		"disclaimer":    r.Disclaimer,
	}
	// 取值必须与 Python 的 RagStatus 一致：pending / ready / degraded / failed
	// （rag-service/app/api/schemas.py）。原先这里写的是 "ok"——那是 Go 侧自造的
	// 第五个值，前端与运维脚本按枚举判断时会漏掉它，看起来像"没有增强结果"。
	if len(d.Citations) == 0 {
		// 一条引用都没过验证时，明确标注"无据"，不要让它看起来像有依据。
		rag["status"] = "degraded"
	} else {
		rag["status"] = "ready"
	}
	return rag
}
