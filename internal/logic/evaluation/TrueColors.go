package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// TrueColors 色彩性格测试评分分析
func (e *Evaluation) TrueColors(data []v1.AnswerRecord) (interface{}, error) {
	// 计算True Colors四个维度的得分
	// 这里简单地将题目分为四个维度，实际应用中需要根据具体的True Colors测试题目分配来计算
	blueScore := 0
	goldScore := 0
	greenScore := 0
	orangeScore := 0

	// 简单的维度分配（实际应用中需要根据具体题目调整）
	for i, record := range data {
		score := int(record.Score)
		switch i % 4 {
		case 0:
			blueScore += score
		case 1:
			goldScore += score
		case 2:
			greenScore += score
		case 3:
			orangeScore += score
		}
	}

	// 确定主导性格类型
	maxScore := blueScore
	dominantType := "蓝色"
	if goldScore > maxScore {
		maxScore = goldScore
		dominantType = "金色"
	}
	if greenScore > maxScore {
		maxScore = greenScore
		dominantType = "绿色"
	}
	if orangeScore > maxScore {
		dominantType = "橙色"
	}

	// 根据主导性格类型提供分析和建议
	level := "完成测试"
	title := dominantType + "性格特质"
	analysis := ""
	var suggestions []string

	switch dominantType {
	case "蓝色":
		analysis = "您的主导性格类型是蓝色，具有强烈的共情力和理想主义倾向。您重视和谐的人际关系，关心他人的感受，富有同情心和同理心。"
		suggestions = []string{
			"学会在关心他人的同时照顾自己的需求",
			"避免过度情绪化，培养理性思考能力",
			"设定明确的个人目标，增强执行力",
			"学会接受生活中的不完美",
		}
	case "金色":
		analysis = "您的主导性格类型是金色，具有强烈的秩序感和责任感。您重视规则和结构，喜欢按计划行事，注重细节和准确性。"
		suggestions = []string{
			"学会灵活变通，接受生活中的变化",
			"培养创新思维，尝试新的方法",
			"学会放松，避免过度追求完美",
			"增强对他人不同观点的包容度",
		}
	case "绿色":
		analysis = "您的主导性格类型是绿色，具有强烈的逻辑思维和创新能力。您重视事实和数据，善于分析问题，喜欢挑战传统思维。"
		suggestions = []string{
			"学会关注他人的感受，增强人际关系能力",
			"避免过度分析，学会直觉判断",
			"培养团队合作精神，学会与他人协作",
			"学会表达自己的情感，增强沟通能力",
		}
	case "橙色":
		analysis = "您的主导性格类型是橙色，具有强烈的活力和适应性。您喜欢尝试新的体验，富有创造力和想象力，善于适应变化。"
		suggestions = []string{
			"学会专注，培养持久力",
			"设定明确的目标和计划，增强执行力",
			"培养责任感，学会承担责任",
			"学会关注细节，避免粗心大意",
		}
	}

	return map[string]interface{}{
		"test_name":    "色彩性格测试（True Colors）",
		"blue_score":   blueScore,
		"gold_score":   goldScore,
		"green_score":  greenScore,
		"orange_score": orangeScore,
		"dominant_type": dominantType,
		"level":        level,
		"title":        title,
		"analysis":     analysis,
		"suggestions":  suggestions,
	}, nil
}