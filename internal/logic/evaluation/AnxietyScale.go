package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// AnxietyScale 焦虑性测试 As量表评分分析
func (e *Evaluation) AnxietyScale(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 计算平均分
	averageScore := float64(totalScore) / float64(len(data))

	// 根据平均分评估焦虑程度
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	switch {
	case averageScore >= 4.5:
		level = "严重焦虑"
		title = "严重焦虑状态"
		analysis = "您的焦虑程度非常高，可能会严重影响您的日常生活和功能。您可能经常感到过度的紧张、担忧和恐惧感，伴随明显的生理症状。"
		suggestions = []string{
			"立即寻求专业的心理健康帮助，如心理咨询或治疗",
			"学习和练习放松技巧，如深呼吸、渐进性肌肉放松等",
			"建立规律的生活作息，保证充足的睡眠",
			"避免咖啡因和其他兴奋剂",
			"与亲友保持联系，寻求社会支持",
		}
	case averageScore >= 3.5:
		level = "中度焦虑"
		title = "中度焦虑状态"
		analysis = "您的焦虑程度中等，可能会对您的日常生活造成一定的影响。您可能经常感到紧张和担忧，伴随一些生理症状。"
		suggestions = []string{
			"寻求专业的心理健康帮助，如心理咨询",
			"学习压力管理技巧，如冥想、瑜伽等",
			"设定合理的目标，避免过度压力",
			"保持健康的生活方式，包括规律的运动和饮食",
			"练习正念和自我觉察，识别焦虑触发因素",
		}
	case averageScore >= 2.5:
		level = "轻度焦虑"
		title = "轻度焦虑状态"
		analysis = "您的焦虑程度较轻，在某些情况下可能会感到紧张和担忧，但通常不会严重影响您的日常生活。"
		suggestions = []string{
			"学习简单的放松技巧，如深呼吸、短暂的冥想等",
			"保持健康的生活方式，包括规律的运动和睡眠",
			"培养积极的应对策略，如问题解决和认知重构",
			"保持良好的时间管理，避免过度承诺",
			"与亲友分享您的感受，寻求支持",
		}
	default:
		level = "低焦虑"
		title = "低焦虑状态"
		analysis = "您的焦虑程度较低，通常能够保持冷静和放松，即使在压力情境下也能有效应对。"
		suggestions = []string{
			"继续保持健康的生活方式和应对策略",
			"定期进行自我评估，关注焦虑水平的变化",
			"学习压力管理技巧，预防未来可能的焦虑增加",
			"培养韧性和心理弹性，增强应对挑战的能力",
			"帮助他人学习有效的焦虑管理技巧",
		}
	}

	return map[string]interface{}{
		"test_name":   "焦虑性测试 As量表",
		"total_score": totalScore,
		"average_score": averageScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}