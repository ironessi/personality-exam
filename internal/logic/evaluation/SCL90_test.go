package evaluation

import (
	"strings"
	"testing"
)

// allOnes 返回 90 题全答“没有”（1 分）的分数切片。
func allOnes() []int {
	scores := make([]int, 90)
	for i := range scores {
		scores[i] = 1
	}
	return scores
}

// TestSCL90CrisisOverridesNormalBand 是这次安全兜底的核心用例。
//
// 「有轻生念头」（第 58 题）答“严重”、其余 89 题全答“没有”时：
// 总分只有 94，总分分档会判成「正常范围」；因子均值也够不到 2.0。
// 危机判定必须把分档结论压过去，否则系统会告诉一个报告了自杀意念的人
// “你在正常范围内”。
func TestSCL90CrisisOverridesNormalBand(t *testing.T) {
	scores := allOnes()
	scores[57] = 5

	result := calculateDetailedResults(scores)

	if result.TotalScore != 94 {
		t.Fatalf("总分 = %d，期望 94", result.TotalScore)
	}
	if !result.NeedFurtherEval {
		t.Error("need_further_eval 应为 true")
	}
	if result.Level != SCL90CrisisLevel {
		t.Errorf("level = %q，期望 %q", result.Level, SCL90CrisisLevel)
	}
	if result.Title == "情绪状态正常" {
		t.Error("title 不应仍是「情绪状态正常」")
	}
	if strings.Contains(result.Analysis, "正常范围") {
		t.Errorf("analysis 不应出现「正常范围」：%s", result.Analysis)
	}
	if !strings.Contains(result.Analysis, "精神科医生") {
		t.Errorf("analysis 缺少求助指引：%s", result.Analysis)
	}

	var found bool
	for _, f := range result.HighRiskFactors {
		if f.Name == SCL90CrisisLabel {
			found = true
		}
	}
	if !found {
		t.Errorf("high_risk_factors 缺少 %q，实得 %v", SCL90CrisisLabel, result.HighRiskFactors)
	}
}

// TestSCL90CrisisThresholdAtItem58 第 58 题的下限是 2：答“很轻”就该触发。
func TestSCL90CrisisThresholdAtItem58(t *testing.T) {
	for score, want := range map[int]bool{1: false, 2: true, 5: true} {
		scores := allOnes()
		scores[57] = score
		got := calculateDetailedResults(scores).Level == SCL90CrisisLevel
		if got != want {
			t.Errorf("第58题答 %d 分：触发 = %v，期望 %v", score, got, want)
		}
	}
}

// TestSCL90CrisisThresholdAtItem20 第 20 题（感到生活无意义）下限是 4。
func TestSCL90CrisisThresholdAtItem20(t *testing.T) {
	for score, want := range map[int]bool{3: false, 4: true} {
		scores := allOnes()
		scores[19] = score
		got := calculateDetailedResults(scores).Level == SCL90CrisisLevel
		if got != want {
			t.Errorf("第20题答 %d 分：触发 = %v，期望 %v", score, got, want)
		}
	}
}

// TestSCL90NoCrisisKeepsNormalBand 确认兜底没有误伤普通作答。
func TestSCL90NoCrisisKeepsNormalBand(t *testing.T) {
	result := calculateDetailedResults(allOnes())
	if result.NeedFurtherEval {
		t.Error("全答“没有”时 need_further_eval 应为 false")
	}
	if result.Level != "正常范围" {
		t.Errorf("level = %q，期望「正常范围」", result.Level)
	}
}
