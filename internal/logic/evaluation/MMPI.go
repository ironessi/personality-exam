package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// MMPI 多相人格测试评分分析
func (e *Evaluation) MMPI(data []v1.AnswerRecord) (interface{}, error) {
	// 计算MMPI各个临床量表的得分
	// 这里简单地将题目分为不同的量表，实际应用中需要根据具体的MMPI测试题目分配来计算
	hypochondriasisScore := 0       // 疑病
	depressionScore := 0            // 抑郁
	hysteriaScore := 0              // 癔症
	psychopathicDeviateScore := 0   // 精神病态
	masculinityFemininityScore := 0 // 男性化-女性化
	paranoiaScore := 0              // 妄想狂
	psychastheniaScore := 0         // 精神衰弱
	schizophreniaScore := 0         // 精神分裂症
	hypomaniaScore := 0             // 轻躁狂
	socialIntroversionScore := 0    // 社会内向性

	// 简单的量表分配（实际应用中需要根据具体题目调整）
	for i, record := range data {
		score := int(record.Score)
		switch i % 10 {
		case 0:
			hypochondriasisScore += score
		case 1:
			depressionScore += score
		case 2:
			hysteriaScore += score
		case 3:
			psychopathicDeviateScore += score
		case 4:
			masculinityFemininityScore += score
		case 5:
			paranoiaScore += score
		case 6:
			psychastheniaScore += score
		case 7:
			schizophreniaScore += score
		case 8:
			hypomaniaScore += score
		case 9:
			socialIntroversionScore += score
		}
	}

	// 计算总分
	totalScore := hypochondriasisScore + depressionScore + hysteriaScore +
		psychopathicDeviateScore + masculinityFemininityScore + paranoiaScore +
		psychastheniaScore + schizophreniaScore + hypomaniaScore + socialIntroversionScore

	// 评估整体心理健康状况
	level := "完成测试"
	title := "MMPI多相人格测试结果"
	analysis := "您的MMPI测试已完成，以下是各个量表的得分情况："
	analysis += "疑病：" + string(rune('0'+hypochondriasisScore)) + "分；"
	analysis += "抑郁：" + string(rune('0'+depressionScore)) + "分；"
	analysis += "癔症：" + string(rune('0'+hysteriaScore)) + "分；"
	analysis += "精神病态：" + string(rune('0'+psychopathicDeviateScore)) + "分；"
	analysis += "男性化-女性化：" + string(rune('0'+masculinityFemininityScore)) + "分；"
	analysis += "妄想狂：" + string(rune('0'+paranoiaScore)) + "分；"
	analysis += "精神衰弱：" + string(rune('0'+psychastheniaScore)) + "分；"
	analysis += "精神分裂症：" + string(rune('0'+schizophreniaScore)) + "分；"
	analysis += "轻躁狂：" + string(rune('0'+hypomaniaScore)) + "分；"
	analysis += "社会内向性：" + string(rune('0'+socialIntroversionScore)) + "分；"
	analysis += "总分：" + string(rune('0'+totalScore)) + "分；"
	analysis += "请注意：MMPI测试结果仅作为参考，不能替代专业的心理诊断。如果您对测试结果有任何疑问，请咨询专业的心理健康专家。"

	// 提供一般性建议
	suggestions := []string{
		"保持健康的生活方式，包括规律的作息和适当的运动",
		"学习有效的压力管理技巧，如冥想、深呼吸等",
		"建立良好的社会支持网络，与亲友保持联系",
		"如果您感到持续的情绪困扰，寻求专业的心理咨询帮助",
		"定期进行自我反思，关注自己的心理健康状况",
	}

	return map[string]interface{}{
		"test_name":                    "多相人格测试",
		"total_score":                  totalScore,
		"hypochondriasis_score":        hypochondriasisScore,
		"depression_score":             depressionScore,
		"hysteria_score":               hysteriaScore,
		"psychopathic_deviate_score":   psychopathicDeviateScore,
		"masculinity_femininity_score": masculinityFemininityScore,
		"paranoia_score":               paranoiaScore,
		"psychasthenia_score":          psychastheniaScore,
		"schizophrenia_score":          schizophreniaScore,
		"hypomania_score":              hypomaniaScore,
		"social_introversion_score":    socialIntroversionScore,
		"level":                        level,
		"title":                        title,
		"analysis":                     analysis,
		"suggestions":                  suggestions,
	}, nil
}
