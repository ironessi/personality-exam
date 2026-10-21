package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"PersonalityTest/internal/consts"
	"encoding/json"
	"testing"
)

// answers 用「相对题号 → 分数」构造作答，base 是这批题目的最小 ID。
// 旧库里题目 ID 是 1..N，新库（public.sql + /api/initialize）里是全局自增的
// （MBTI 是 376..385）。
func answers(base int, scores []int64) []v1.AnswerRecord {
	out := make([]v1.AnswerRecord, 0, len(scores))
	for i, score := range scores {
		out = append(out, v1.AnswerRecord{
			QuestionID: base + i,
			Answer:     "A",
			Score:      score,
		})
	}
	return out
}

func mbtiScores() []int64 {
	return []int64{5, 2, 6, 5, 5, 2, 5, 4, 5, 2} // E15 S6 T11 J9
}

// TestMBTIScoringIndependentOfQuestionIDBase 是这个回归的核心：
// 评分曾经按硬编码的题目 ID（1/5/9、2/6/10…）分维度，题目 ID 一旦不是
// 从 1 开始（新库是 376..385），四个维度全部为 0，任何作答都判成 INFP。
func TestMBTIScoringIndependentOfQuestionIDBase(t *testing.T) {
	e := New()

	for _, base := range []int{1, 376, 5001} {
		got, err := e.Evaluate(consts.MBTI, answers(base, mbtiScores()))
		if err != nil {
			t.Fatalf("base=%d: Evaluate 返回错误: %v", base, err)
		}
		m, ok := got.(map[string]interface{})
		if !ok {
			t.Fatalf("base=%d: 返回值类型不是 map[string]interface{}", base)
		}

		want := map[string]int{
			"extraversion": 15,
			"sensing":      6,
			"thinking":     11,
			"judging":      9,
			"total_score":  41,
		}
		for key, wantVal := range want {
			gotVal, _ := m[key].(int)
			if gotVal != wantVal {
				t.Errorf("base=%d: %s = %d，期望 %d", base, key, gotVal, wantVal)
			}
		}
		if typ, _ := m["personality_type"].(string); typ != "ENTJ" {
			t.Errorf("base=%d: personality_type = %q，期望 ENTJ", base, typ)
		}
	}
}

// TestEvaluateNormalizesQuestionIDsForAllEvaluators 覆盖其余按「题号 1..N」
// 写的评分函数：Personalities 写死 case 1..20、PhilPersonality 写死 case 1..10、
// RIASEC 用 questionID%6、AMS 用 questionID%2、SBTI 用 map[1..32]。
// 判据是「题目 ID 整体平移不改变任何输出」——比逐个断言维度分数更能抓住
// 六维旋转、奇偶翻转这类不会报错的错算。
func TestEvaluateNormalizesQuestionIDsForAllEvaluators(t *testing.T) {
	cases := []struct {
		name  string
		code  int
		count int
	}{
		{"MBTI", consts.MBTI, 10},
		{"Personalities", consts.Personalities, 20},
		{"PhilPersonality", consts.PhilPersonality, 10},
		{"RIASEC", consts.RIASEC, 108},
		{"AMS", consts.AMS, 10},
		{"SBTI", consts.SBTI, 32},
	}

	e := New()
	for _, tc := range cases {
		scores := make([]int64, tc.count)
		for i := range scores {
			scores[i] = int64(1 + i%4)
		}

		base1, err := e.Evaluate(tc.code, answers(1, scores))
		if err != nil {
			t.Fatalf("%s: base=1 评分失败: %v", tc.name, err)
		}
		shifted, err := e.Evaluate(tc.code, answers(376, scores))
		if err != nil {
			t.Fatalf("%s: base=376 评分失败: %v", tc.name, err)
		}

		wantJSON, _ := json.Marshal(base1)
		gotJSON, _ := json.Marshal(shifted)
		if string(wantJSON) != string(gotJSON) {
			t.Errorf("%s: 题目 ID 从 376 起算后结果不同\n base=1  : %s\n base=376: %s",
				tc.name, wantJSON, gotJSON)
		}
	}
}

// TestNormalizeQuestionIDsKeepsOneBasedInput 保证旧库（题号从 1 开始）不被
// 额外拷贝或改写，避免给老数据引入行为变化。
func TestNormalizeQuestionIDsKeepsOneBasedInput(t *testing.T) {
	in := answers(1, []int64{3, 4})
	out := normalizeQuestionIDs(in)

	if &in[0] != &out[0] {
		t.Error("从 1 开始时应原样返回同一个切片，不做拷贝")
	}
	if out[0].QuestionID != 1 || out[1].QuestionID != 2 {
		t.Errorf("题目 ID 被改写: %d, %d", out[0].QuestionID, out[1].QuestionID)
	}
}

// TestNormalizeQuestionIDsDoesNotMutateInput 保证归一化不会改到调用方的切片
// （controller 之后还要把原始作答写进 user_answers）。
func TestNormalizeQuestionIDsDoesNotMutateInput(t *testing.T) {
	in := answers(376, []int64{3, 4})
	out := normalizeQuestionIDs(in)

	if in[0].QuestionID != 376 || in[1].QuestionID != 377 {
		t.Errorf("原始作答被就地修改: %d, %d", in[0].QuestionID, in[1].QuestionID)
	}
	if out[0].QuestionID != 1 || out[1].QuestionID != 2 {
		t.Errorf("归一化结果不对: %d, %d", out[0].QuestionID, out[1].QuestionID)
	}
}
