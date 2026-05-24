package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// BDI 贝克抑郁自评量表评分分析
func (e *Evaluation) BDI(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 贝克抑郁自评量表评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	if totalScore < 10 {
		level = "无抑郁"
		title = "情绪状态良好"
		analysis = "恭喜！您的情绪状态良好，没有明显的抑郁症状。"
		suggestions = []string{
			"继续保持积极的生活态度和健康的生活方式",
			"定期进行自我情绪评估，保持心理健康",
			"与亲友保持良好的沟通和联系",
		}
	} else if totalScore < 20 {
		level = "轻度抑郁"
		title = "轻度抑郁倾向"
		analysis = "您可能存在轻度抑郁倾向，需要关注自己的情绪变化。"
		suggestions = []string{
			"尝试增加户外活动和运动，有助于改善情绪",
			"培养兴趣爱好，丰富生活内容",
			"与信任的人分享您的感受",
			"保持规律的作息时间",
		}
	} else if totalScore < 30 {
		level = "中度抑郁"
		title = "中度抑郁倾向"
		analysis = "您可能存在中度抑郁倾向，建议寻求专业帮助。"
		suggestions = []string{
			"考虑咨询心理医生或心理咨询师",
			"与家人朋友沟通，获取支持和理解",
			"避免独处，多参与社交活动",
			"保持健康的饮食和睡眠习惯",
		}
	} else {
		level = "重度抑郁"
		title = "重度抑郁倾向"
		analysis = "您可能存在重度抑郁倾向，强烈建议立即寻求专业帮助。"
		suggestions = []string{
			"尽快联系心理医生或精神科医生",
			"告知家人或朋友您的情况，寻求陪伴和支持",
			"避免做出重大决策或单独行动",
			"记住：抑郁是可以治疗的，您并不孤单",
		}
	}

	return map[string]interface{}{
		"test_name":   "贝克抑郁自评量表(BDI)",
		"total_score": totalScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
