package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// Machiavellianism 马基雅维利人格测试评分分析
func (e *Evaluation) Machiavellianism(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 马基雅维利人格测试评分标准
	maxScore := len(data) * 5
	percentage := float64(totalScore) / float64(maxScore) * 100

	if percentage >= 80 {
		level := "高度马基雅维利主义"
		title := "马基雅维利主义倾向高"
		analysis := "您的马基雅维利主义倾向较高，可能更注重结果而非道德，善于操纵和利用他人。"
		suggestions := []string{
			"注意自己的行为对他人的影响，避免过度操纵和利用他人",
			"培养道德意识和同理心，关注他人的感受和需求",
			"学会建立真诚的人际关系，而非仅仅利用他人",
			"考虑寻求专业心理咨询师的帮助，探索自己的行为动机",
		}
		return map[string]interface{}{
			"test_name":   "马基雅维利人格测试",
			"total_score": totalScore,
			"max_score":   maxScore,
			"percentage":  percentage,
			"level":       level,
			"title":       title,
			"analysis":    analysis,
			"suggestions": suggestions,
		}, nil
	} else if percentage >= 60 {
		level := "中等马基雅维利主义"
		title := "马基雅维利主义倾向中等"
		analysis := "您的马基雅维利主义倾向处于中等水平，在某些情况下可能会表现出实用主义和操纵性。"
		suggestions := []string{
			"保持自我觉察，注意自己在人际关系中的行为和动机",
			"平衡结果导向和道德原则，避免过度功利主义",
			"培养合作精神，学会与他人共赢",
			"关注长期关系的建立，而非短期利益",
		}
		return map[string]interface{}{
			"test_name":   "马基雅维利人格测试",
			"total_score": totalScore,
			"max_score":   maxScore,
			"percentage":  percentage,
			"level":       level,
			"title":       title,
			"analysis":    analysis,
			"suggestions": suggestions,
		}, nil
	} else {
		level := "低度马基雅维利主义"
		title := "马基雅维利主义倾向低"
		analysis := "您的马基雅维利主义倾向较低，更注重道德原则和真诚的人际关系。"
		suggestions := []string{
			"继续保持真诚和道德的行为方式",
			"学会保护自己，避免被他人操纵和利用",
			"在适当情况下可以更加果断和自信",
			"平衡理想主义和现实需求，提高自己的适应性",
		}
		return map[string]interface{}{
			"test_name":   "马基雅维利人格测试",
			"total_score": totalScore,
			"max_score":   maxScore,
			"percentage":  percentage,
			"level":       level,
			"title":       title,
			"analysis":    analysis,
			"suggestions": suggestions,
		}, nil
	}
}
