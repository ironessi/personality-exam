package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// ObedienceTest 服从性测试评分分析
func (e *Evaluation) ObedienceTest(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 计算平均分
	averageScore := float64(totalScore) / float64(len(data))

	// 根据平均分评估服从倾向程度
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	switch {
	case averageScore >= 4.5:
		level = "高度服从"
		title = "高度服从倾向"
		analysis = "您的服从倾向程度很高，在面对权威和规则时表现出强烈的服从意愿，可能会影响您的自主性和批判性思维能力。"
		suggestions = []string{
			"学习独立思考，培养批判性思维能力",
			"学会在适当的情况下质疑权威",
			"设定个人边界，学会说'不'",
			"培养自我价值感，相信自己的判断",
			"寻求专业帮助，如心理咨询",
		}
	case averageScore >= 3.5:
		level = "中度服从"
		title = "中度服从倾向"
		analysis = "您的服从倾向程度中等，在大多数情况下会服从权威和规则，但在某些情况下能够表达自己的意见。"
		suggestions = []string{
			"在安全的环境中练习表达自己的意见",
			"学习识别不合理的要求",
			"培养决策能力，增强自主性",
			"学会平衡服从与个人判断",
		}
	case averageScore >= 2.5:
		level = "轻度服从"
		title = "轻度服从倾向"
		analysis = "您的服从倾向程度较轻，在面对权威和规则时能够保持一定的自主性，同时尊重合理的要求。"
		suggestions = []string{
			"继续保持独立思考的能力",
			"学会在团队中有效沟通",
			"培养合作精神，同时保持个人特色",
			"学会尊重他人的意见，同时坚持自己的原则",
		}
	default:
		level = "低服从倾向"
		title = "低服从倾向"
		analysis = "您的服从倾向程度较低，具有较强的自主性和独立性，可能在团队合作中表现出挑战权威的倾向。"
		suggestions = []string{
			"学会在适当的情况下服从合理的权威",
			"培养团队合作精神，尊重集体决策",
			"学会倾听他人的意见，避免过度自我中心",
			"在表达不同意见时保持尊重和建设性",
		}
	}

	return map[string]interface{}{
		"test_name":   "服从性测试",
		"total_score": totalScore,
		"average_score": averageScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}