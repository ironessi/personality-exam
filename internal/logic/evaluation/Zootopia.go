package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// Zootopia 疯狂动物城测试性格评分分析
func (e *Evaluation) Zootopia(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 疯狂动物城测试性格评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	if totalScore <= 15 {
		level = "外向冒险型"
		title = "像朱迪一样勇敢"
		analysis = "您的性格类似于《疯狂动物城》中的朱迪警官，勇敢、乐观、充满正义感，喜欢挑战新事物。"
		suggestions = []string{
			"继续保持积极向上的态度，勇敢追求自己的目标",
			"学会团队合作，发挥领导力优势",
			"在冒险中注意安全，合理评估风险",
		}
	} else if totalScore <= 25 {
		level = "平衡协调型"
		title = "像尼克一样灵活"
		analysis = "您的性格类似于《疯狂动物城》中的尼克，聪明、灵活、善于交际，能够适应各种环境。"
		suggestions = []string{
			"发挥您的社交能力，建立良好的人际关系网络",
			"保持创新思维，寻找解决问题的新方法",
			"培养责任感，将聪明才智用于有意义的事情",
		}
	} else {
		level = "内向思考型"
		title = "像闪电一样沉稳"
		analysis = "您的性格类似于《疯狂动物城》中的闪电，沉稳、专注、有条理，注重细节和质量。"
		suggestions = []string{
			"继续发挥您的专注力和细致性优势",
			"适当增加社交活动，拓展人际关系",
			"在保证质量的同时，提高工作效率",
		}
	}

	return map[string]interface{}{
		"test_name":   "疯狂动物城测试性格",
		"total_score": totalScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
