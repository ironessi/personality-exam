package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// EPPS 爱德华个人偏好测验评分分析
func (e *Evaluation) EPPS(data []v1.AnswerRecord) (interface{}, error) {
	// 计算15种需求的得分
	achievementScore := 0
	conformanceScore := 0
	orderScore := 0
	exhibitionScore := 0
	autonomyScore := 0
	affiliationScore := 0
	intraceptionScore := 0
	succoranceScore := 0
	dominanceScore := 0
	abasementScore := 0
	nurturanceScore := 0
	changeScore := 0
	persistenceScore := 0
	heterosexualityScore := 0
	aggressionScore := 0

	// 简单的需求分配（实际应用中需要根据具体题目调整）
	for i, record := range data {
		score := int(record.Score)
		switch i % 15 {
		case 0:
			achievementScore += score
		case 1:
			conformanceScore += score
		case 2:
			orderScore += score
		case 3:
			exhibitionScore += score
		case 4:
			autonomyScore += score
		case 5:
			affiliationScore += score
		case 6:
			intraceptionScore += score
		case 7:
			succoranceScore += score
		case 8:
			dominanceScore += score
		case 9:
			abasementScore += score
		case 10:
			nurturanceScore += score
		case 11:
			changeScore += score
		case 12:
			persistenceScore += score
		case 13:
			heterosexualityScore += score
		case 14:
			aggressionScore += score
		}
	}

	// 确定主导需求
	maxScore := achievementScore
	dominantNeed := "成就"
	if conformanceScore > maxScore {
		maxScore = conformanceScore
		dominantNeed = "顺从"
	}
	if orderScore > maxScore {
		maxScore = orderScore
		dominantNeed = "秩序"
	}
	if exhibitionScore > maxScore {
		maxScore = exhibitionScore
		dominantNeed = "表现"
	}
	if autonomyScore > maxScore {
		maxScore = autonomyScore
		dominantNeed = "自主"
	}
	if affiliationScore > maxScore {
		maxScore = affiliationScore
		dominantNeed = "亲和"
	}
	if intraceptionScore > maxScore {
		maxScore = intraceptionScore
		dominantNeed = "省察"
	}
	if succoranceScore > maxScore {
		maxScore = succoranceScore
		dominantNeed = "求助"
	}
	if dominanceScore > maxScore {
		maxScore = dominanceScore
		dominantNeed = "支配"
	}
	if abasementScore > maxScore {
		maxScore = abasementScore
		dominantNeed = "谦卑"
	}
	if nurturanceScore > maxScore {
		maxScore = nurturanceScore
		dominantNeed = "慈善"
	}
	if changeScore > maxScore {
		maxScore = changeScore
		dominantNeed = "变通"
	}
	if persistenceScore > maxScore {
		maxScore = persistenceScore
		dominantNeed = "坚毅"
	}
	if heterosexualityScore > maxScore {
		maxScore = heterosexualityScore
		dominantNeed = "异性恋"
	}
	if aggressionScore > maxScore {
		maxScore = aggressionScore
		dominantNeed = "攻击"
	}

	return map[string]interface{}{
		"test_name":             "爱德华个人偏好测验",
		"achievement_score":     achievementScore,
		"conformance_score":     conformanceScore,
		"order_score":           orderScore,
		"exhibition_score":      exhibitionScore,
		"autonomy_score":        autonomyScore,
		"affiliation_score":     affiliationScore,
		"intraception_score":    intraceptionScore,
		"succorance_score":      succoranceScore,
		"dominance_score":       dominanceScore,
		"abasement_score":       abasementScore,
		"nurturance_score":      nurturanceScore,
		"change_score":          changeScore,
		"persistence_score":     persistenceScore,
		"heterosexuality_score": heterosexualityScore,
		"aggression_score":      aggressionScore,
		"dominant_need":         dominantNeed,
		"level":                 "完成测试",
		"title":                 dominantNeed + "需求主导",
		"analysis":              "您的爱德华个人偏好测验已完成，主导需求为" + dominantNeed + "。",
		"suggestions": []string{
			"了解自己的需求模式有助于个人成长",
			"在满足自己需求的同时尊重他人需求",
			"寻求专业心理咨询以深入探索需求动机",
		},
	}, nil
}
