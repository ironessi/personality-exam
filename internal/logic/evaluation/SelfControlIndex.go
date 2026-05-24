package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// SelfControlIndex 自我控力指数评估评分分析
func (e *Evaluation) SelfControlIndex(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 计算平均分
	averageScore := float64(totalScore) / float64(len(data))

	// 根据平均分评估自控力水平
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	switch {
	case averageScore >= 4.5:
		level = "优秀"
		title = "自控力优秀"
		analysis = "您的自控力水平非常高，能够很好地控制自己的行为、情绪和注意力，为了长期目标而牺牲即时满足，并且能够在压力下保持冷静。"
		suggestions = []string{
			"继续保持您的自控力优势",
			"将您的自控力应用到更具挑战性的目标上",
			"帮助他人培养自控力",
			"注意避免过度自我控制导致的压力",
		}
	case averageScore >= 3.5:
		level = "良好"
		title = "自控力良好"
		analysis = "您的自控力水平良好，能够在大多数情况下控制自己的行为和情绪，但在某些情况下可能会遇到挑战。"
		suggestions = []string{
			"在薄弱环节加强自控力训练",
			"学习更有效的压力管理技巧",
			"设定明确的目标和计划",
			"培养健康的习惯来支持自控力",
		}
	case averageScore >= 2.5:
		level = "一般"
		title = "自控力一般"
		analysis = "您的自控力水平一般，在一些情况下能够控制自己，但在面对诱惑或压力时可能会遇到困难。"
		suggestions = []string{
			"从小目标开始，逐步培养自控力",
			"学习识别和避免触发失控的情境",
			"寻求支持和 accountability伙伴",
			"练习 mindfulness 和自我觉察",
		}
	default:
		level = "需要提高"
		title = "自控力需要提高"
		analysis = "您的自控力水平较低，在面对诱惑、压力或困难时容易失去控制，可能会影响您的目标实现和生活质量。"
		suggestions = []string{
			"寻求专业帮助，如心理咨询",
			"学习基础的自控力训练技巧",
			"创建支持性的环境，减少诱惑",
			"设定具体、可实现的短期目标",
			"培养自我同情，避免因挫折而放弃",
		}
	}

	return map[string]interface{}{
		"test_name":   "自我控力指数评估",
		"total_score": totalScore,
		"average_score": averageScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}