package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// StanfordBinet 斯坦福-比奈智力量表评分分析
func (e *Evaluation) StanfordBinet(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 斯坦福-比奈智力量表评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	if totalScore >= 36 {
		level = "优秀"
		title = "智力水平优秀"
		analysis = "您的智力水平优秀，具有出色的逻辑思维、语言理解和问题解决能力。"
		suggestions = []string{
			"挑战自己，设定更高的学习和职业目标",
			"培养创新思维，尝试解决复杂问题",
			"保持终身学习的习惯，不断拓展知识领域",
		}
	} else if totalScore >= 24 {
		level = "良好"
		title = "智力水平良好"
		analysis = "您的智力水平良好，具有较强的逻辑思维和问题解决能力。"
		suggestions = []string{
			"继续加强学习，提高知识储备",
			"培养批判性思维能力，多角度思考问题",
			"参加智力挑战活动，保持大脑活跃",
		}
	} else if totalScore >= 12 {
		level = "中等"
		title = "智力水平中等"
		analysis = "您的智力水平中等，具有基本的逻辑思维和问题解决能力。"
		suggestions = []string{
			"注重基础知识的学习和巩固",
			"通过练习提高逻辑思维能力",
			"保持积极的学习态度，相信自己的潜力",
		}
	} else {
		level = "待提高"
		title = "智力水平待提高"
		analysis = "您的智力水平有待提高，建议加强学习和思维训练。"
		suggestions = []string{
			"从基础开始，循序渐进地学习",
			"多做逻辑思维和问题解决的练习",
			"保持耐心和毅力，相信通过努力可以提高",
		}
	}

	return map[string]interface{}{
		"test_name":   "斯坦福-比奈智力量表",
		"total_score": totalScore,
		"max_score":   40,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
