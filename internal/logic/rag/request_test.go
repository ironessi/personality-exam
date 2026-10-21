package rag

import (
	"encoding/json"
	"strings"
	"testing"

	"PersonalityTest/internal/consts"
)

// mbtiResult 是 MBTI 评分函数真实返回形状的最小副本（internal/logic/evaluation/MBTI.go）。
func mbtiResult() map[string]interface{} {
	return map[string]interface{}{
		"test_name":        "16型人格测试（MBTI）",
		"extraversion":     15,
		"sensing":          6,
		"thinking":         11,
		"judging":          9,
		"total_score":      41,
		"personality_type": "ENTJ",
		"level":            "完成测试",
	}
}

// scl90Result 是 SCL-90 评分函数真实返回形状的最小副本。
// 关键点：factors 是**以中文因子名为键的 map**（SCL90.go 里是
// factorDetails[factor.Name] = {...}），不是切片。早期夹具写成切片，
// 导致真实提交在 Prepare 阶段找不到因子，而单测全绿。
func scl90Result() map[string]interface{} {
	names := []string{
		"躯体化", "强迫症状", "人际关系敏感", "抑郁", "焦虑",
		"敌对", "恐怖", "偏执", "精神病性", "其他",
	}
	scores := []float64{2.4, 1.8, 1.6, 3.1, 2.2, 1.2, 1.4, 1.9, 1.3, 1.7}
	factors := make(map[string]interface{}, len(names))
	for i, name := range names {
		factors[name] = map[string]interface{}{
			"score": scores[i], "raw_score": 10, "items": []int{1, 2},
		}
	}
	return map[string]interface{}{
		"test_name":         "SCL-90症状自评量表",
		"total_score":       198,
		"max_score":         450,
		"factors":           factors,
		"level":             "轻度异常",
		"need_further_eval": true,
	}
}

func TestBuildRequestMBTI(t *testing.T) {
	req, err := BuildRequest(BuildInput{
		TestType: consts.MBTI, TestName: "16型人格测试（MBTI）", Result: mbtiResult(),
	})
	if err != nil {
		t.Fatalf("BuildRequest 失败: %v", err)
	}

	if req.SchemaVersion != 1 {
		t.Errorf("SchemaVersion = %d，期望 1", req.SchemaVersion)
	}
	if req.Test.Code != consts.MBTI || req.Test.Key != "mbti" {
		t.Errorf("test = %+v，期望 code=13 key=mbti", req.Test)
	}
	if req.Test.Policy != "full" {
		t.Errorf("policy = %q，期望 full", req.Test.Policy)
	}
	// 检索范围 = 本域 + common，顺序固定。
	if len(req.Test.Domains) != 2 || req.Test.Domains[0] != "typology" || req.Test.Domains[1] != "common" {
		t.Errorf("domains = %v，期望 [typology common]", req.Test.Domains)
	}
	if len(req.Measures) != 4 {
		t.Fatalf("measures 数量 = %d，期望 4", len(req.Measures))
	}

	// strength 口径：(score - threshold) / (max - threshold)，与
	// rag-service/scripts/e2e_enhance.py 的 _strength_audit 一致。
	want := []struct {
		key      string
		score    float64
		max      float64
		thr      float64
		pole     string
		strength float64
	}{
		{"EI", 15, 21, 12, "E", 0.3333},
		{"SN", 6, 21, 12, "N", -0.6667},
		{"TF", 11, 14, 8, "T", 0.5},
		{"JP", 9, 14, 8, "J", 0.1667},
	}
	for i, w := range want {
		m := req.Measures[i]
		if m.Key != w.key {
			t.Errorf("measures[%d].Key = %q，期望 %q", i, m.Key, w.key)
		}
		if m.Score != w.score {
			t.Errorf("%s.Score = %v，期望 %v", w.key, m.Score, w.score)
		}
		if m.Max == nil || *m.Max != w.max {
			t.Errorf("%s.Max = %v，期望 %v", w.key, m.Max, w.max)
		}
		if m.Threshold == nil || *m.Threshold != w.thr {
			t.Errorf("%s.Threshold = %v，期望 %v", w.key, m.Threshold, w.thr)
		}
		if m.Pole != w.pole {
			t.Errorf("%s.Pole = %q，期望 %q", w.key, m.Pole, w.pole)
		}
		if m.Strength == nil || absDiff(*m.Strength, w.strength) > 0.001 {
			t.Errorf("%s.Strength = %v，期望 %v", w.key, m.Strength, w.strength)
		}
		// factor 专用字段不得出现在 dimension 上。
		if m.ClinicalCutoff != nil || m.FactorLevel != "" {
			t.Errorf("%s 是 dimension，不该有 clinical_cutoff/factor_level", w.key)
		}
	}

	if req.Labels["type"] != "ENTJ" {
		t.Errorf("labels.type = %q，期望 ENTJ", req.Labels["type"])
	}
	if len(req.Flags) != 0 {
		t.Errorf("flags = %v，MBTI 无风险信号，期望空", req.Flags)
	}

	// MBTI 结果里有 total_score（41），要如实带上；没有 max_score 就留空。
	if req.Result.Total == nil || *req.Result.Total != 41 {
		t.Errorf("result.total = %v，期望 41", req.Result.Total)
	}
	if req.Result.TotalMax != nil {
		t.Errorf("result.total_max = %v，MBTI 没有上界，期望 nil", req.Result.TotalMax)
	}
}

func TestBuildRequestSCL90(t *testing.T) {
	req, err := BuildRequest(BuildInput{
		TestType: consts.SCL90, TestName: "SCL-90症状自评量表", Result: scl90Result(),
	})
	if err != nil {
		t.Fatalf("BuildRequest 失败: %v", err)
	}

	if len(req.Measures) != 10 {
		t.Fatalf("measures 数量 = %d，期望 10", len(req.Measures))
	}
	for _, m := range req.Measures {
		if m.Kind != "factor" {
			t.Errorf("%s.Kind = %q，期望 factor", m.Key, m.Kind)
		}
		if m.ClinicalCutoff == nil || *m.ClinicalCutoff != 2.0 {
			t.Errorf("%s.ClinicalCutoff = %v，期望 2.0", m.Key, m.ClinicalCutoff)
		}
		// factor 不应带维度字段。
		if m.Threshold != nil || m.Max != nil || m.Pole != "" {
			t.Errorf("%s 是 factor，不该有 threshold/max/pole", m.Key)
		}
	}

	// 抑郁 3.1 分：strength = (3.1-2)/2 = 0.55，且判为偏高。
	dep := findMeasure(t, req.Measures, "depression")
	if dep.Score != 3.1 {
		t.Errorf("depression.Score = %v，期望 3.1（按中文因子名匹配 factors[]）", dep.Score)
	}
	if dep.FactorLevel != "偏高" {
		t.Errorf("depression.FactorLevel = %q，期望 偏高", dep.FactorLevel)
	}
	if dep.Strength == nil || absDiff(*dep.Strength, 0.55) > 0.001 {
		t.Errorf("depression.Strength = %v，期望 0.55", dep.Strength)
	}

	// 敌对 1.2 < 2.0：低于临床参考线。
	hos := findMeasure(t, req.Measures, "hostility")
	if hos.FactorLevel != "正常" {
		t.Errorf("hostility.FactorLevel = %q，期望 正常", hos.FactorLevel)
	}
	if hos.Strength == nil || absDiff(*hos.Strength, -0.4) > 0.001 {
		t.Errorf("hostility.Strength = %v，期望 -0.4", hos.Strength)
	}

	// need_further_eval 必须转成 needs_referral，否则危机提示不会触发检索。
	if len(req.Flags) != 1 || req.Flags[0] != "needs_referral" {
		t.Errorf("flags = %v，期望 [needs_referral]", req.Flags)
	}

	// SCL-90 有 max_score=450，要带上。
	if req.Result.TotalMax == nil || *req.Result.TotalMax != 450 {
		t.Errorf("result.total_max = %v，期望 450", req.Result.TotalMax)
	}
}

// TestMeasureStrengthIsClamped 钉死 [-1,1] 的夹取：SCL-90 满分 5 分时
// (5-2)/2 = 1.5，必须夹成 1.0，否则 Python 侧 ge/le 校验直接 422。
func TestMeasureStrengthIsClamped(t *testing.T) {
	result := mbtiResult()
	result["extraversion"] = 21 // (21-12)/9 = 1.0
	result["sensing"] = 0       // (0-12)/9 = -1.333 → 夹成 -1.0
	req, err := BuildRequest(BuildInput{TestType: consts.MBTI, Result: result})
	if err != nil {
		t.Fatalf("BuildRequest 失败: %v", err)
	}
	if s := *findMeasure(t, req.Measures, "EI").Strength; s != 1.0 {
		t.Errorf("EI.Strength = %v，期望 1.0", s)
	}
	if s := *findMeasure(t, req.Measures, "SN").Strength; s != -1.0 {
		t.Errorf("SN.Strength = %v，期望 -1.0", s)
	}
}

// TestBuildRequestRejectsTestsWithoutMeasures 钉死「没有阈值就不增强」：
// DISC 只比四维谁最高，没有绝对阈值，宁可跳过也不编一个分档点。
func TestBuildRequestRejectsTestsWithoutMeasures(t *testing.T) {
	_, err := BuildRequest(BuildInput{
		TestType: consts.DISC, Result: map[string]interface{}{"d_score": 10},
	})
	if err == nil {
		t.Fatal("DISC 没有测量项定义，BuildRequest 应该报错")
	}
	if !strings.Contains(err.Error(), "测量项") {
		t.Errorf("错误信息应说明缺测量项，实际: %v", err)
	}
}

// TestBuildRequestRejectsDisabledPolicy 钉死 disabled 策略不走增强。
func TestBuildRequestRejectsDisabledPolicy(t *testing.T) {
	_, err := BuildRequest(BuildInput{
		TestType: consts.Zootopia, Result: map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("Zootopia 的 policy 是 disabled，BuildRequest 应该报错")
	}
}

// TestBuildRequestFailsOnMissingField：评分代码改了字段名而注册表没同步时，
// 必须**报错**而不是静默发一个 score=0——那会让报告说「你的外倾维度是 0 分」。
func TestBuildRequestFailsOnMissingField(t *testing.T) {
	result := mbtiResult()
	delete(result, "extraversion")
	_, err := BuildRequest(BuildInput{TestType: consts.MBTI, Result: result})
	if err == nil {
		t.Fatal("缺 extraversion 时应该报错")
	}
	if !strings.Contains(err.Error(), "extraversion") {
		t.Errorf("错误信息应包含缺失的字段名，实际: %v", err)
	}
}

// TestBuildRequestJSONShape 检查序列化后**没有多余字段**：
// Python 侧 Measure / EnhanceRequest 都是 extra="forbid"，
// 多一个键就是 422，而不是被忽略。
func TestBuildRequestJSONShape(t *testing.T) {
	req, err := BuildRequest(BuildInput{TestType: consts.MBTI, Result: mbtiResult()})
	if err != nil {
		t.Fatalf("BuildRequest 失败: %v", err)
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	// 顶层只允许这 7 个键。
	allowedTop := map[string]bool{
		"schema_version": true, "test": true, "result": true,
		"measures": true, "labels": true, "flags": true, "options": true,
	}
	for k := range decoded {
		if !allowedTop[k] {
			t.Errorf("请求体出现契约外的顶层键 %q", k)
		}
	}

	// dimension 的 JSON 里不得出现 factor 专用键（omitempty 应已省掉）。
	for _, m := range decoded["measures"].([]interface{}) {
		mm := m.(map[string]interface{})
		for _, forbidden := range []string{"clinical_cutoff", "factor_level"} {
			if _, ok := mm[forbidden]; ok {
				t.Errorf("dimension %v 的 JSON 里出现了 %q", mm["key"], forbidden)
			}
		}
	}
}

func TestDomainsOfCommonOnly(t *testing.T) {
	// common_only 的测评只能检索 common：本域语料就算存在也不该被召回，
	// 否则等于变相为受保护测验做个性化解读（ADR-0008）。
	meta := consts.TestRegistry[consts.Rorschach]
	got := domainsOf(meta)
	if len(got) != 1 || got[0] != "common" {
		t.Errorf("domainsOf(Rorschach) = %v，期望 [common]", got)
	}
}

func findMeasure(t *testing.T, list []Measure, key string) Measure {
	t.Helper()
	for _, m := range list {
		if m.Key == key {
			return m
		}
	}
	t.Fatalf("找不到测量项 %q", key)
	return Measure{}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
