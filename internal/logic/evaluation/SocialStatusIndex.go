package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

// SocialStatusIndex 社会地位指数测试ST量表评分分析
func (e *Evaluation) SocialStatusIndex(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	// 社会地位指数测试ST量表评分标准
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	maxScore := len(data) * 5
	percentage := float64(totalScore) / float64(maxScore) * 100

	if percentage >= 80 {
		level = "高社会地位"
		title = "社会地位指数高"
		analysis = "恭喜！您的社会地位指数较高，在经济、职业和社会资本方面都有较好的表现。"
		suggestions = []string{
			"继续保持和提升自己的社会地位",
			"利用自己的社会资源帮助他人，回馈社会",
			"关注社会公平和公益事业，承担更多的社会责任",
		}
	} else if percentage >= 60 {
		level = "中等社会地位"
		title = "社会地位指数中等"
		analysis = "您的社会地位指数处于中等水平，在经济、职业和社会资本方面有一定的积累。"
		suggestions = []string{
			"继续努力提升自己的教育水平和职业技能",
			"拓展自己的社会关系网络，建立更多的社会资本",
			"关注职业发展，寻求更高的职业地位和收入水平",
		}
	} else if percentage >= 40 {
		level = "较低社会地位"
		title = "社会地位指数较低"
		analysis = "您的社会地位指数较低，需要关注自己在经济、职业和社会资本方面的提升。"
		suggestions = []string{
			"努力提高自己的教育水平，获取更多的知识和技能",
			"关注就业机会，寻找适合自己的职业发展方向",
			"积极参与社会活动，拓展自己的社会关系网络",
			"学习理财知识，合理规划自己的经济生活",
		}
	} else {
		level = "低社会地位"
		title = "社会地位指数低"
		analysis = "您的社会地位指数较低，建议尽快采取行动提升自己的社会地位。"
		suggestions = []string{
			"寻求职业培训和教育机会，提升自己的就业能力",
			"咨询职业规划师，寻找适合自己的职业发展道路",
			"加入社区组织或志愿者团体，拓展自己的社会网络",
			"学习基本的理财知识，改善自己的经济状况",
		}
	}

	return map[string]interface{}{
		"test_name":   "社会地位指数测试ST量表",
		"total_score": totalScore,
		"max_score":   maxScore,
		"percentage":  percentage,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}