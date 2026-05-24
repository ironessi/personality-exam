package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// MHLAQ 心理健康素养量表评分分析
func (e *Evaluation) MHLAQ(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 心理健康素养量表评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	maxScore := len(data) * 5
	percentage := float64(totalScore) / float64(maxScore) * 100

	if percentage >= 80 {
		level = "高水平"
		title = "心理健康素养水平高"
		analysis = "恭喜！您的心理健康素养水平较高，对心理健康知识有较好的了解和应用能力。"
		suggestions = []string{
			"继续保持对心理健康知识的学习和应用",
			"分享您的心理健康知识，帮助他人提升心理健康素养",
			"关注最新的心理健康研究和信息",
		}
	} else if percentage >= 60 {
		level = "中等水平"
		title = "心理健康素养水平中等"
		analysis = "您的心理健康素养水平处于中等，对心理健康知识有一定的了解和应用能力。"
		suggestions = []string{
			"继续学习心理健康知识，提升自己的心理健康素养",
			"尝试将所学的心理健康知识应用到日常生活中",
			"关注自己的心理健康，定期进行自我评估",
		}
	} else if percentage >= 40 {
		level = "较低水平"
		title = "心理健康素养水平较低"
		analysis = "您的心理健康素养水平较低，建议加强对心理健康知识的学习。"
		suggestions = []string{
			"参加心理健康知识培训或课程",
			"阅读心理健康相关书籍和文章",
			"关注专业心理健康机构的公众号或网站",
			"与心理健康专业人士交流，获取指导和建议",
		}
	} else {
		level = "低水平"
		title = "心理健康素养水平低"
		analysis = "您的心理健康素养水平较低，建议尽快加强对心理健康知识的学习。"
		suggestions = []string{
			"寻求专业心理健康教育和培训",
			"与心理健康专业人士建立联系，获取个性化的指导",
			"优先学习基础的心理健康知识，如情绪管理、压力应对等",
			"加入心理健康支持小组，与他人分享学习经验",
		}
	}

	return map[string]interface{}{
		"test_name":   "心理健康素养量表(MHLAQ)",
		"total_score": totalScore,
		"max_score":   maxScore,
		"percentage":  percentage,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}