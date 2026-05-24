package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// NPI 自恋人格量表评分分析
func (e *Evaluation) NPI(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 计算平均分
	averageScore := float64(totalScore) / float64(len(data))

	// 根据平均分评估自恋程度
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	switch {
	case averageScore >= 4.5:
		level = "高度自恋"
		title = "高度自恋倾向"
		analysis = "您的自恋倾向程度非常高，可能会表现出强烈的自我中心、优越感和对他人的缺乏同理心。这可能会影响您的人际关系和社会功能。"
		suggestions = []string{
			"寻求专业心理咨询帮助",
			"学习培养同理心和关注他人的感受",
			"练习倾听和尊重他人的意见",
			"培养团队合作精神",
			"学习接受批评和反馈",
		}
	case averageScore >= 3.5:
		level = "中度自恋"
		title = "中度自恋倾向"
		analysis = "您的自恋倾向程度中等，在某些情况下表现出自我中心和优越感，但在其他情况下能够关注他人的需求。"
		suggestions = []string{
			"在关注自己的同时也关注他人的需求",
			"学习平等地与他人交流",
			"培养感恩和谦逊的态度",
			"尝试从他人的角度看问题",
		}
	case averageScore >= 2.5:
		level = "轻度自恋"
		title = "轻度自恋倾向"
		analysis = "您的自恋倾向程度较轻，在大多数情况下能够平衡自我关注和对他人的关注，但在某些特定情况下可能会表现出自恋特征。"
		suggestions = []string{
			"保持健康的自我价值感",
			"继续培养对他人的同理心",
			"学会欣赏他人的成就",
			"平衡自我表达和倾听他人",
		}
	default:
		level = "低自恋"
		title = "低自恋倾向"
		analysis = "您的自恋倾向程度较低，通常能够关注他人的需求，表现出谦逊和同理心，重视人际关系的平衡。"
		suggestions = []string{
			"继续保持对他人的同理心和关注",
			"培养健康的自我价值感",
			"在适当的情况下表达自己的需求和成就",
			"学会为自己的成就感到自豪",
		}
	}

	return map[string]interface{}{
		"test_name":   "自恋人格量表（NPI）",
		"total_score": totalScore,
		"average_score": averageScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}