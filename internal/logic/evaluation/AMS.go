package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// AMSResult 成就动机测试(AMS)结果
type AMSResult struct {
	HopeForSuccess int      // 希望成功的动机(MS)
	FearOfFailure  int      // 害怕失败的动机(MF)
	TotalScore     int      // 总分
	Level          string   // 等级
	Title          string   // 标题
	Analysis       string   // 分析
	Suggestions    []string // 建议
}

// AMS 成就动机测试(AMS)评分分析
func (e *Evaluation) AMS(data []v1.AnswerRecord) (interface{}, error) {
	// 计算希望成功的动机(MS)和害怕失败的动机(MF)得分
	hopeForSuccess := 0
	fearOfFailure := 0

	// 根据题目分配维度
	// 奇数序号题目: 希望成功的动机(MS)
	// 偶数序号题目: 害怕失败的动机(MF)
	for _, record := range data {
		score := int(record.Score)
		if record.QuestionID%2 == 1 {
			// 奇数题目: 希望成功的动机
			hopeForSuccess += score
		} else {
			// 偶数题目: 害怕失败的动机
			fearOfFailure += score
		}
	}

	totalScore := hopeForSuccess + fearOfFailure

	// 确定等级和分析
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	title = "成就动机测试(AMS)结果"
	analysis = fmt.Sprintf("您的成就动机测试得分情况如下：\n希望成功的动机(MS)：%d分\n害怕失败的动机(MF)：%d分\n总分：%d分\n\n各维度解释：\n- 希望成功的动机(MS)：反映个体追求成功、完成具有挑战性任务的倾向\n- 害怕失败的动机(MF)：反映个体避免失败、减少负面评价的倾向\n\n成就动机类型：", 
		hopeForSuccess, fearOfFailure, totalScore)

	// 确定成就动机类型
	if hopeForSuccess > fearOfFailure {
		analysis += "\n您属于高希望成功型：倾向于追求挑战，主动设定高目标，享受成功的喜悦"
	} else if fearOfFailure > hopeForSuccess {
		analysis += "\n您属于高害怕失败型：倾向于避免挑战，害怕失败带来的负面评价，追求安全感"
	} else {
		analysis += "\n您属于平衡型：希望成功和害怕失败的动机相当，能灵活应对不同情境"
	}

	// 提供建议
	suggestions = []string{
		"了解自己的成就动机类型有助于制定更合适的目标",
		"高希望成功型可以尝试更多挑战性任务，发挥自己的潜力",
		"高害怕失败型可以逐步设定小目标，积累成功经验",
		"平衡型可以根据情境灵活调整策略，保持良好的心理状态",
		"定期进行成就动机测试，追踪自己的动机变化",
	}

	return map[string]interface{}{
		"test_name":       "成就动机测试(AMS)",
		"hope_for_success": hopeForSuccess,
		"fear_of_failure":  fearOfFailure,
		"total_score":     totalScore,
		"level":           level,
		"title":           title,
		"analysis":        analysis,
		"suggestions":     suggestions,
	}, nil
}