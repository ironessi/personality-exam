package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// EMP 同理心通情量表评分分析
func (e *Evaluation) EMP(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 同理心通情量表评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	maxScore := len(data) * 5
	percentage := float64(totalScore) / float64(maxScore) * 100

	if percentage >= 80 {
		level = "高度同理心"
		title = "同理心水平高"
		analysis = "恭喜！您的同理心水平较高，能够很好地理解和分享他人的情感状态。"
		suggestions = []string{
			"继续发挥您的同理心优势，帮助他人建立良好的人际关系",
			"考虑从事需要高同理心的职业，如心理咨询师、教师、医护人员等",
			"学习如何保护自己的情感边界，避免情感耗竭",
		}
	} else if percentage >= 60 {
		level = "中等同理心"
		title = "同理心水平中等"
		analysis = "您的同理心水平处于中等，能够较好地理解和分享他人的情感状态。"
		suggestions = []string{
			"继续培养同理心，尝试更多地倾听和理解他人",
			"学习非语言沟通技巧，提高情感感知能力",
			"参加同理心训练课程或工作坊，进一步提升同理心水平",
		}
	} else if percentage >= 40 {
		level = "较低同理心"
		title = "同理心水平较低"
		analysis = "您的同理心水平较低，需要关注自己的同理心培养。"
		suggestions = []string{
			"尝试站在他人的角度思考问题，培养换位思考的能力",
			"学习情感识别和表达技巧，提高情感认知能力",
			"参加同理心训练活动，如角色扮演、小组讨论等",
			"阅读关于同理心的书籍，了解同理心的重要性和培养方法",
		}
	} else {
		level = "低同理心"
		title = "同理心水平低"
		analysis = "您的同理心水平较低，建议尽快关注自己的同理心培养。"
		suggestions = []string{
			"寻求专业心理咨询师的帮助，探索同理心发展的障碍",
			"参加结构化的同理心训练课程，系统提升同理心水平",
			"从简单的倾听开始，逐步培养理解他人的能力",
			"关注他人的情感需求，尝试表达对他人的关心和理解",
		}
	}

	return map[string]interface{}{
		"test_name":   "同理心通情量表（EMP）",
		"total_score": totalScore,
		"max_score":   maxScore,
		"percentage":  percentage,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
