package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// HSP 敏感度测试评分分析
func (e *Evaluation) HSP(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 计算平均分
	averageScore := float64(totalScore) / float64(len(data))

	// 根据平均分评估敏感度程度
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	switch {
	case averageScore >= 4.5:
		level = "高度敏感"
		title = "高度敏感人群（HSP）"
		analysis = "您的敏感度程度非常高，属于高度敏感人群（HSP）。您对环境刺激、情绪和感官信息的处理深度和反应强度明显高于平均水平。"
		suggestions = []string{
			"了解并接受自己的敏感性，这是您的优势",
			"创建安静、舒适的生活环境",
			"学习设置个人边界，避免过度刺激",
			"练习正念和自我照顾技巧",
			"寻找理解和支持您敏感性的社交圈",
		}
	case averageScore >= 3.5:
		level = "中度敏感"
		title = "中度敏感人群"
		analysis = "您的敏感度程度中等，在某些情况下表现出较高的敏感性，但在其他情况下能够适应。"
		suggestions = []string{
			"识别并尊重自己的敏感触发因素",
			"在需要时给自己创造安静的空间",
			"学习有效的压力管理技巧",
			"平衡社交活动和独处时间",
		}
	case averageScore >= 2.5:
		level = "轻度敏感"
		title = "轻度敏感人群"
		analysis = "您的敏感度程度较轻，在大多数情况下能够适应环境刺激，但在某些特定情况下可能会表现出敏感反应。"
		suggestions = []string{
			"关注自己的身体和情绪信号",
			"在压力大时给自己适当的休息",
			"培养健康的生活方式",
			"学会欣赏自己的敏感性带来的优势",
		}
	default:
		level = "低敏感"
		title = "低敏感人群"
		analysis = "您的敏感度程度较低，通常能够适应各种环境刺激，对情绪和感官信息的反应强度相对较低。"
		suggestions = []string{
			"继续保持您的适应能力",
			"尝试理解他人的敏感需求",
			"培养同理心，关注他人的情绪状态",
			"在需要时给自己一些安静的时间",
		}
	}

	return map[string]interface{}{
		"test_name":   "敏感度测试（HSP）",
		"total_score": totalScore,
		"average_score": averageScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}