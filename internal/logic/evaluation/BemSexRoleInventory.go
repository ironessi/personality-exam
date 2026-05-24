package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// BemSexRoleInventory 贝姆性别角色测试量表评分分析
func (e *Evaluation) BemSexRoleInventory(data []v1.AnswerRecord) (interface{}, error) {
	// 计算男性化和女性化特质的得分
	masculinityScore := 0
	femininityScore := 0
	androgynyScore := 0

	// 简单的特质分配（实际应用中需要根据具体题目调整）
	for i, record := range data {
		score := int(record.Score)
		switch i % 3 {
		case 0:
			masculinityScore += score
		case 1:
			femininityScore += score
		case 2:
			androgynyScore += score
		}
	}

	// 确定性别角色类型
	genderRoleType := ""
	if masculinityScore > femininityScore {
		if masculinityScore-femininityScore > 10 {
			genderRoleType = "男性化"
		} else {
			genderRoleType = "双性化"
		}
	} else if femininityScore > masculinityScore {
		if femininityScore-masculinityScore > 10 {
			genderRoleType = "女性化"
		} else {
			genderRoleType = "双性化"
		}
	} else {
		genderRoleType = "双性化"
	}

	// 提供分析和建议
	level := "完成测试"
	title := genderRoleType + "性别角色特质"
	analysis := "您的贝姆性别角色测试已完成，以下是您的得分情况：男性化特质：" + string(rune('0'+masculinityScore)) + "分；女性化特质：" + string(rune('0'+femininityScore)) + "分；双性化特质：" + string(rune('0'+androgynyScore)) + "分；您的性别角色类型为：" + genderRoleType

	var suggestions []string
	switch genderRoleType {
	case "男性化":
		suggestions = []string{
			"尝试培养更多的女性化特质，如同理心和情感表达",
			"学习倾听和理解他人的感受",
			"在竞争的同时也注重合作",
			"培养更多的耐心和细致",
		}
	case "女性化":
		suggestions = []string{
			"尝试培养更多的男性化特质，如自信和果断",
			"学习表达自己的意见和需求",
			"在合作的同时也培养竞争意识",
			"培养更多的独立性和自主性",
		}
	case "双性化":
		suggestions = []string{
			"继续保持平衡的性别角色特质",
			"在不同情境中灵活运用不同的特质",
			"欣赏自己的多样性和适应性",
			"帮助他人打破性别刻板印象",
		}
	}

	return map[string]interface{}{
		"test_name":         "贝姆性别角色测试量表",
		"masculinity_score": masculinityScore,
		"femininity_score":  femininityScore,
		"androgyny_score":   androgynyScore,
		"gender_role_type":  genderRoleType,
		"level":             level,
		"title":             title,
		"analysis":          analysis,
		"suggestions":       suggestions,
	}, nil
}
