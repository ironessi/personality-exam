package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// PassivePersonality 被动型人格测试评分分析
func (e *Evaluation) PassivePersonality(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 计算平均分
	averageScore := float64(totalScore) / float64(len(data))

	// 根据平均分评估被动倾向程度
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	switch {
	case averageScore >= 4.5:
		level = "高度被动"
		title = "高度被动型人格"
		analysis = "您的被动倾向程度很高，在决策、人际互动和责任承担等场景中表现出明显的被动特征，可能会影响您的个人发展和人际关系。"
		suggestions = []string{
			"学习识别和表达自己的需求和感受",
			"练习做出小决策，逐步培养决策能力",
			"学习设定个人边界，学会说'不'",
			"寻求专业帮助，如心理咨询",
			"参加自信心训练课程",
		}
	case averageScore >= 3.5:
		level = "中度被动"
		title = "中度被动型人格"
		analysis = "您的被动倾向程度中等，在某些场景中表现出被动特征，但在其他场景中能够主动应对。"
		suggestions = []string{
			"在安全的环境中练习表达自己的意见",
			"学习识别被动行为模式并尝试改变",
			"培养自我觉察能力，了解自己的价值观和需求",
			"设定明确的个人目标，增强主动性",
		}
	case averageScore >= 2.5:
		level = "轻度被动"
		title = "轻度被动型人格"
		analysis = "您的被动倾向程度较轻，在大多数场景中能够主动应对，但在某些特定情况下可能会表现出被动特征。"
		suggestions = []string{
			"识别触发被动行为的特定情境",
			"在这些情境中练习主动应对策略",
			"培养自信，相信自己的判断和能力",
			"与支持性的人建立关系，获得反馈和鼓励",
		}
	default:
		level = "低被动倾向"
		title = "低被动倾向人格"
		analysis = "您的被动倾向程度较低，在大多数场景中能够主动应对，表现出较强的自主性和决策能力。"
		suggestions = []string{
			"继续保持您的主动性和自主性",
			"学习倾听他人的意见，避免过度主导",
			"在团队合作中平衡个人主动性和集体需求",
			"帮助他人培养主动性和自信心",
		}
	}

	return map[string]interface{}{
		"test_name":   "被动型人格测试",
		"total_score": totalScore,
		"average_score": averageScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}