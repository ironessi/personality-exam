package rag

import (
	"fmt"
	"math"

	"PersonalityTest/internal/consts"
)

// 本文件把「Go 侧的测评结果」翻译成 Python 契约里的 EnhanceRequest。
//
// 三件必须守住的事：
//  1. **只发契约里有的键**。rag-service 的 EnhanceRequest / Measure 都是
//     extra="forbid"，多发一个字段就是 422，而不是被忽略。
//  2. **strength 由 Go 算**。Python 侧不解析量纲、不判断测评类型；它只信
//     strength（控制流程）与 factor_level（喂 prompt）。把这两个的算法放在
//     Go 是因为阈值和满分只有在 Go 的评分代码里才有权威值。
//  3. **绝不编造阈值**。没有阈值定义的测评（如 DISC）宁可跳过增强，也不
//     随便给一个分档点——那会变成伪造的"专业判断"。

// schemaVersion 必须与 rag-service/app/api/schemas.py 的
// SUPPORTED_SCHEMA_VERSION 一致。契约变更时要显式升级，不要静默兼容。
const schemaVersion = 1

// BuildInput 是拼请求所需的全部输入。
type BuildInput struct {
	TestType       int
	TestName       string                 // 优先用 DB 里的 test_types.name
	Result         map[string]interface{} // evaluation 返回并补全后的 resultMap
	MaxSuggestions int                    // 0 表示用默认值
}

// EnhanceRequest 是 POST /api/v1/enhance 的请求体。
//
// 字段名与顺序刻意与 schemas.py 一一对应，便于对照排查 422。
type EnhanceRequest struct {
	SchemaVersion int               `json:"schema_version"`
	Test          TestInfo          `json:"test"`
	Result        ResultInfo        `json:"result"`
	Measures      []Measure         `json:"measures"`
	Labels        map[string]string `json:"labels"`
	Flags         []string          `json:"flags"`
	Options       EnhanceOptions    `json:"options"`
}

// TestInfo 描述"这是什么测评、该检索哪些域"。
//
// domains 是**构念域**（ADR-0008），不是测评名。policy=common_only 时
// 必须恰好是 ["common"]，policy=full 时必须包含 "common"——这是 Python
// 侧 TestInfo 校验器的硬要求。
type TestInfo struct {
	Code    int      `json:"code"`
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	Policy  string   `json:"policy"`
}

// ResultInfo 是整体结果，仅用于叙述。total 的量纲随量表而异，
// **不参与算强度**（强度只看 measures[].strength）。
type ResultInfo struct {
	Level    string   `json:"level"`
	Total    *float64 `json:"total,omitempty"`
	TotalMax *float64 `json:"total_max,omitempty"`
}

// Measure 统一承载 dimension 与 factor，用 Kind 判别。
//
// 指针字段用于"这个量纲不该出现"的场景：factor 不发 threshold/pole，
// dimension 不发 clinical_cutoff。虽然 Python 侧接受 null，但省掉更干净。
type Measure struct {
	Key      string   `json:"key"`
	Label    string   `json:"label"`
	Kind     string   `json:"kind"`
	Score    float64  `json:"score"`
	Max      *float64 `json:"max,omitempty"`
	Pole     string   `json:"pole,omitempty"`
	Strength *float64 `json:"strength,omitempty"`

	// dimension 专用
	Threshold *float64 `json:"threshold,omitempty"`

	// factor 专用
	ClinicalCutoff *float64 `json:"clinical_cutoff,omitempty"`
	FactorLevel    string   `json:"factor_level,omitempty"`
}

// EnhanceOptions 是检索与生成的选项。
type EnhanceOptions struct {
	TopK             int    `json:"top_k"`
	MaxSuggestions   int    `json:"max_suggestions"`
	Language         string `json:"language"`
	IncludeCitations bool   `json:"include_citations"`
}

// 默认值：top_k=5 与 rag-service 的 Settings.top_k 一致；
// max_suggestions 对应 config.yaml 的 rag.maxSuggestions。
const (
	defaultTopK = 5
	// fallbackMaxSuggestions 只在调用方没传 MaxSuggestions 时兜底；
	// 正常路径下它来自 rag.maxSuggestions（config.go 的 defaultMaxSuggestions）。
	fallbackMaxSuggestions = 5
)

// BuildRequest 拼出规范请求。
//
// 返回错误的情形只有两类，都是**配置/注册表问题**而不是用户输入问题：
//   - 该测评没有增强元数据（未登记或 disabled）；
//   - 该测评没有测量项定义（目前是 DISC，见 consts.TestMeasures 注释）。
//
// 这两类都应该被 service.go 的闸门提前拦掉，走到这里报错说明闸门漏了。
func BuildRequest(in BuildInput) (*EnhanceRequest, error) {
	meta, ok := consts.LookupTest(in.TestType)
	if !ok {
		return nil, fmt.Errorf("测评 %d 没有登记增强元数据（ADR-0008）", in.TestType)
	}
	if meta.Policy == consts.PolicyDisabled {
		return nil, fmt.Errorf("测评 %d 的增强策略是 disabled", in.TestType)
	}

	specs := consts.MeasuresOf(in.TestType)
	if len(specs) == 0 {
		return nil, fmt.Errorf(
			"测评 %d（%s）没有测量项定义，无法拼出 measures[]；"+
				"Python 侧 measures 有 min_length=1，且 dimension 强制带 threshold",
			in.TestType, meta.Key,
		)
	}

	measures := make([]Measure, 0, len(specs))
	for _, spec := range specs {
		m, err := buildMeasure(spec, in.Result)
		if err != nil {
			return nil, err
		}
		measures = append(measures, m)
	}

	name := in.TestName
	if name == "" {
		name = meta.Name
	}

	maxSuggestions := in.MaxSuggestions
	if maxSuggestions <= 0 {
		maxSuggestions = fallbackMaxSuggestions
	}

	return &EnhanceRequest{
		SchemaVersion: schemaVersion,
		Test: TestInfo{
			Code:    meta.Code,
			Key:     meta.Key,
			Name:    name,
			Domains: domainsOf(meta),
			Policy:  string(meta.Policy),
		},
		Result:   buildResultInfo(in.Result),
		Measures: measures,
		Labels:   buildLabels(meta, in.Result),
		Flags:    buildFlags(meta, in.Result),
		Options: EnhanceOptions{
			TopK:             defaultTopK,
			MaxSuggestions:   maxSuggestions,
			Language:         "zh-CN",
			IncludeCitations: true,
		},
	}, nil
}

// domainsOf 把构念域元数据翻成契约要的字符串数组。
//
// ⚠️ common_only 的测评（罗夏、斯坦福-比奈）**不能**带本域：它们的
// 本域语料就算存在也不该被检索，因为那会诱导模型对该测验做个性化解读。
func domainsOf(meta consts.TestMeta) []string {
	if meta.Policy == consts.PolicyCommonOnly {
		return []string{string(consts.DomainCommon)}
	}
	domains := []string{string(meta.Domain)}
	if meta.Domain != consts.DomainCommon {
		domains = append(domains, string(consts.DomainCommon))
	}
	return domains
}

// buildMeasure 从 resultMap 里取一个测量项并算好 strength。
func buildMeasure(spec consts.MeasureSpec, result map[string]interface{}) (Measure, error) {
	score, ok := lookupScore(spec, result)
	if !ok {
		return Measure{}, fmt.Errorf(
			"测评结果里找不到测量项 %q 的取值（Source=%q, FromFactors=%v）；"+
				"这通常意味着评分代码改了字段名而注册表没同步",
			spec.Key, spec.Source, spec.FromFactors,
		)
	}

	m := Measure{
		Key:   spec.Key,
		Label: spec.Label,
		Kind:  string(spec.Kind),
		Score: score,
	}

	switch spec.Kind {
	case consts.MeasureDimension:
		if spec.Max <= spec.Threshold {
			return Measure{}, fmt.Errorf(
				"测量项 %q 的 Max(%.2f) 必须大于 Threshold(%.2f)",
				spec.Key, spec.Max, spec.Threshold,
			)
		}
		threshold := spec.Threshold
		max := spec.Max
		m.Threshold = &threshold
		m.Max = &max
		if score >= spec.Threshold {
			m.Pole = spec.HighPole
		} else {
			m.Pole = spec.LowPole
		}
		// 口径与 rag-service/scripts/e2e_enhance.py 的 _strength_audit 一致：
		//   strength = (score - threshold) / (max - threshold)，再 clamp 到 [-1,1]
		// 用「距分档点的差 / 分档点以上的跨度」而不是 score/max：
		// 后者会让"刚好过线"看起来只有 0.5 强，把两侧的差别抹平。
		m.Strength = clampUnit((score - spec.Threshold) / (spec.Max - spec.Threshold))

	case consts.MeasureFactor:
		if spec.ClinicalCutoff <= 0 {
			return Measure{}, fmt.Errorf("测量项 %q 的 ClinicalCutoff 必须为正", spec.Key)
		}
		cutoff := spec.ClinicalCutoff
		m.ClinicalCutoff = &cutoff
		m.Strength = clampUnit((score - cutoff) / cutoff)
		if score >= cutoff {
			m.FactorLevel = "偏高"
		} else {
			m.FactorLevel = "正常"
		}

	default:
		return Measure{}, fmt.Errorf("测量项 %q 的 Kind %q 不合法", spec.Key, spec.Kind)
	}

	return m, nil
}

// lookupScore 按注册表声明的位置取分。
func lookupScore(spec consts.MeasureSpec, result map[string]interface{}) (float64, bool) {
	if !spec.FromFactors {
		return toFloat(result[spec.Source])
	}

	// factors 有两种真实形态，都必须认：
	//   1. map 以因子名为键 —— 这是 SCL90.go 实际返回的形态：
	//      factors["躯体化"] = {score, raw_score, items}
	//   2. 切片，元素带 name 字段 —— 早期单测夹具用的形态。
	// 只认第 2 种会让真实 SCL-90 提交在 Prepare 阶段就失败（找不到
	// somatization），而单测因为是手写夹具反而通过。
	raw, ok := result["factors"]
	if !ok {
		return 0, false
	}
	if byName, ok := raw.(map[string]interface{}); ok {
		item, ok := byName[spec.Source]
		if !ok {
			return 0, false
		}
		if m, ok := item.(map[string]interface{}); ok {
			return toFloat(m["score"])
		}
		return 0, false
	}
	list, ok := raw.([]map[string]interface{})
	if !ok {
		// 经过一次 JSON 往返后元素类型会变成 map[string]interface{} 的切片，
		// 但如果上游用的是 []interface{} 也要认。
		generic, ok := raw.([]interface{})
		if !ok {
			return 0, false
		}
		list = make([]map[string]interface{}, 0, len(generic))
		for _, item := range generic {
			if m, ok := item.(map[string]interface{}); ok {
				list = append(list, m)
			}
		}
	}
	for _, item := range list {
		if name, _ := item["name"].(string); name == spec.Source {
			return toFloat(item["score"])
		}
	}
	return 0, false
}

// buildResultInfo 取整体结果。level 是所有 41 个评分函数都返回的五个共有键之一。
func buildResultInfo(result map[string]interface{}) ResultInfo {
	info := ResultInfo{}
	if level, ok := result["level"].(string); ok {
		info.Level = level
	}
	if total, ok := toFloat(result["total_score"]); ok {
		info.Total = &total
	}
	// max_score 只有部分测评有（SCL-90 写死 450）。没有就不发，让 Python
	// 用 None——**不要自己推一个**，量纲不同推出错的比没有更糟。
	if totalMax, ok := toFloat(result["max_score"]); ok {
		info.TotalMax = &totalMax
	}
	return info
}

// buildLabels 收集结论性的标签（如 MBTI 的 INTJ），供 prompt 叙述使用。
//
// ⚠️ 这些标签**不进 measures**，因此不参与 strength 计算，也不会被当作
// 检索 query。它们只是让模型知道"这是个 ENFP"。
func buildLabels(meta consts.TestMeta, result map[string]interface{}) map[string]string {
	labels := make(map[string]string, 2)
	if t, ok := result["personality_type"].(string); ok && t != "" {
		labels["type"] = t
	}
	if meta.Policy == consts.PolicyCommonOnly {
		// 这两个测验不做个性化解读，标签只用于"避免误判"提示。
		labels["scope"] = "common_only"
	}
	return labels
}

// buildFlags 收集安全标记。
//
// needs_referral 有两个来源：
//   - 注册表强制（common_only 的两个受保护测验）；
//   - 评分结果本身（SCL-90 的 need_further_eval、命中的危机信号）。
//
// Python 侧看到 needs_referral 会额外检索「何时需要寻求专业心理帮助」，
// 所以这个标记**宁可多给不可漏给**。
func buildFlags(meta consts.TestMeta, result map[string]interface{}) []string {
	flags := make([]string, 0, 1)
	if consts.NeedsReferralForced(meta.Code) {
		flags = append(flags, "needs_referral")
		return flags
	}
	if need, ok := result["need_further_eval"].(bool); ok && need {
		flags = append(flags, "needs_referral")
	}
	return flags
}

// toFloat 把 interface{} 统一成 float64。
//
// 为什么需要它：resultMap 里的数字来源有三类——评分代码直接塞的 int、
// float64，以及经过 json.Unmarshal 之后的 float64。同一个字段在不同路径上
// 可能是 int 也可能是 float64，写死类型断言会静默取到 0。
func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	case nil:
		return 0, false
	default:
		return 0, false
	}
}

// clampUnit 把值夹到 [-1,1]。
func clampUnit(v float64) *float64 {
	if math.IsNaN(v) {
		return nil
	}
	if v > 1 {
		v = 1
	} else if v < -1 {
		v = -1
	}
	return &v
}
