package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// Rorschach 罗夏墨迹测验评分分析
func (e *Evaluation) Rorschach(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 罗夏墨迹测验评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	if totalScore <= 15 {
		level = "外向情感型"
		title = "情感表达丰富"
		analysis = "您的性格倾向于外向，情感表达丰富，善于社交，对外部世界充满兴趣。"
		suggestions = []string{
			"发挥您的社交优势，建立广泛的人际关系",
			"学会情绪管理，保持情感平衡",
			"在社交活动中注意倾听他人，避免过度自我表达",
		}
	} else if totalScore <= 25 {
		level = "平衡综合型"
		title = "心理状态平衡"
		analysis = "您的性格较为平衡，能够较好地适应不同情境，具有一定的自我洞察力。"
		suggestions = []string{
			"继续保持心理平衡，培养情绪调节能力",
			"加强自我反思，提高自我认知水平",
			"在面对挑战时保持冷静，理性分析问题",
		}
	} else {
		level = "内向思考型"
		title = "思维深度丰富"
		analysis = "您的性格倾向于内向，思维深度丰富，善于思考和分析，对内部世界更为关注。"
		suggestions = []string{
			"适当增加社交活动，避免过度自我封闭",
			"学会表达自己的想法和感受，与他人分享内心世界",
			"保持思维的开放性，接纳不同的观点和想法",
		}
	}

	return map[string]interface{}{
		"test_name":   "罗夏墨迹测验（Rorschach）",
		"total_score": totalScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
