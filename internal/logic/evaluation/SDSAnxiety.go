package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// SDSAnxiety 焦虑症自测量表评分分析
func (e *Evaluation) SDSAnxiety(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 焦虑症自测量表评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	if totalScore < 20 {
		level = "无焦虑"
		title = "情绪状态稳定"
		analysis = "恭喜！您的情绪状态稳定，没有明显的焦虑症状。"
		suggestions = []string{
			"继续保持积极的生活态度和健康的生活方式",
			"定期进行自我情绪评估，保持心理健康",
			"与亲友保持良好的沟通和联系",
		}
	} else if totalScore < 40 {
		level = "轻度焦虑"
		title = "轻度焦虑倾向"
		analysis = "您可能存在轻度焦虑倾向，需要关注自己的情绪变化。"
		suggestions = []string{
			"尝试学习一些放松技巧，如深呼吸、冥想等",
			"保持规律的作息时间，确保充足的睡眠",
			"适当增加户外活动和运动，有助于缓解焦虑",
			"避免过度使用咖啡因和酒精",
		}
	} else if totalScore < 60 {
		level = "中度焦虑"
		title = "中度焦虑倾向"
		analysis = "您可能存在中度焦虑倾向，建议寻求专业帮助。"
		suggestions = []string{
			"考虑咨询心理医生或心理咨询师",
			"与家人朋友沟通，获取支持和理解",
			"尝试认知行为疗法，改变负面思维模式",
			"保持健康的饮食和睡眠习惯",
		}
	} else {
		level = "重度焦虑"
		title = "重度焦虑倾向"
		analysis = "您可能存在重度焦虑倾向，强烈建议立即寻求专业帮助。"
		suggestions = []string{
			"尽快联系心理医生或精神科医生",
			"告知家人或朋友您的情况，寻求陪伴和支持",
			"避免做出重大决策或单独行动",
			"记住：焦虑是可以治疗的，您并不孤单",
		}
	}

	return map[string]interface{}{
		"test_name":   "焦虑症自测量表(SDS)",
		"total_score": totalScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
