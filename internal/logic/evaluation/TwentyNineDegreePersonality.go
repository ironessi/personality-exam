package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// TwentyNineDegreePersonality 29度人格测试评分分析
func (e *Evaluation) TwentyNineDegreePersonality(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 29度人格测试评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	maxScore := len(data) * 5
	percentage := float64(totalScore) / float64(maxScore) * 100

	if percentage >= 80 {
		level = "高成熟度"
		title = "人格特质成熟度高"
		analysis = "恭喜！您的人格特质成熟度较高，能够很好地平衡终结性与开创性的张力。"
		suggestions = []string{
			"继续发挥您的人格优势，在各个领域取得更大的成就",
			"尝试挑战更高难度的任务，进一步提升自己",
			"分享您的成长经验，帮助他人成长",
		}
	} else if percentage >= 60 {
		level = "中等成熟度"
		title = "人格特质成熟度中等"
		analysis = "您的人格特质成熟度处于中等水平，能够较好地应对生活中的各种挑战。"
		suggestions = []string{
			"关注自己的人格成长，尝试突破固有模式",
			"学习新的技能和知识，拓展自己的视野",
			"培养良好的情绪管理能力，提高应对压力的能力",
		}
	} else if percentage >= 40 {
		level = "较低成熟度"
		title = "人格特质成熟度较低"
		analysis = "您的人格特质成熟度较低，需要关注自己的人格成长和发展。"
		suggestions = []string{
			"尝试了解自己的人格特质，认识自己的优势和不足",
			"学习人格心理学相关知识，提升自我认知",
			"寻求专业心理咨询师的帮助，探索人格成长的方向",
			"参加人格成长培训或工作坊，加速成长进程",
		}
	} else {
		level = "低成熟度"
		title = "人格特质成熟度低"
		analysis = "您的人格特质成熟度较低，建议尽快关注自己的人格成长和发展。"
		suggestions = []string{
			"立即寻求专业心理咨询师的帮助，进行深入的人格探索",
			"阅读人格成长相关书籍，建立正确的自我认知",
			"参加结构化的人格成长课程，系统提升人格成熟度",
			"与信任的人分享您的感受和困惑，获取支持和理解",
		}
	}

	return map[string]interface{}{
		"test_name":   "29度人格测试",
		"total_score": totalScore,
		"max_score":   maxScore,
		"percentage":  percentage,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}