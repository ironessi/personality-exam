package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// HCL32 轻躁狂症状自评量表评分分析
func (e *Evaluation) HCL32(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 轻躁狂症状自评量表评分标准
	if totalScore < 14 {
		level := "无轻躁狂倾向"
		title := "情绪状态稳定"
		analysis := "恭喜！您的情绪状态稳定，没有明显的轻躁狂症状。"
		suggestions := []string{
			"继续保持积极的生活态度和健康的生活方式",
			"定期进行自我情绪评估，保持心理健康",
			"与亲友保持良好的沟通和联系",
		}
		return map[string]interface{}{
			"test_name":   "轻躁狂症状自评量表(HCL-32)",
			"total_score": totalScore,
			"level":       level,
			"title":       title,
			"analysis":    analysis,
			"suggestions": suggestions,
		}, nil
	} else if totalScore < 21 {
		level := "轻度轻躁狂倾向"
		title := "轻度轻躁狂倾向"
		analysis := "您可能存在轻度轻躁狂倾向，需要关注自己的情绪变化。"
		suggestions := []string{
			"注意观察自己的情绪波动，尤其是兴奋、冲动等症状",
			"保持规律的作息时间，确保充足的睡眠",
			"避免过度劳累和压力过大",
			"与家人朋友沟通，获取支持和理解",
		}
		return map[string]interface{}{
			"test_name":   "轻躁狂症状自评量表(HCL-32)",
			"total_score": totalScore,
			"level":       level,
			"title":       title,
			"analysis":    analysis,
			"suggestions": suggestions,
		}, nil
	} else {
		level := "中度至重度轻躁狂倾向"
		title := "中度至重度轻躁狂倾向"
		analysis := "您可能存在中度至重度轻躁狂倾向，建议寻求专业帮助。"
		suggestions := []string{
			"尽快联系心理医生或精神科医生进行评估",
			"告知家人或朋友您的情况，寻求陪伴和支持",
			"避免做出重大决策或冒险行为",
			"保持健康的生活方式，包括规律的饮食和睡眠",
		}
		return map[string]interface{}{
			"test_name":   "轻躁狂症状自评量表(HCL-32)",
			"total_score": totalScore,
			"level":       level,
			"title":       title,
			"analysis":    analysis,
			"suggestions": suggestions,
		}, nil
	}
}
