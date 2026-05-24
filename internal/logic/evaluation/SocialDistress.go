package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) SocialDistress(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// 社交回避及苦恼量表分级标准
	if totalScore <= 8 {
		level = "低度社交困扰"
		title = "社交状态良好"
		analysis = fmt.Sprintf("您的社交回避及苦恼得分为%d分（总分28分），表明您在社交场合中相对自在，没有明显的社交困扰。", totalScore)
		suggestions = []string{
			"继续保持健康的社交习惯",
			"适度参与社交活动",
		}
	} else if totalScore <= 16 {
		level = "轻度社交困扰"
		title = "社交状态基本良好"
		analysis = fmt.Sprintf("您的社交回避及苦恼得分为%d分（总分28分），表明您在某些社交场合会感到不适，但总体可控。", totalScore)
		suggestions = []string{
			"尝试逐步扩大舒适区",
			"学习放松技巧应对社交压力",
			"在安全的环境中练习社交技能",
		}
	} else if totalScore <= 22 {
		level = "中度社交困扰"
		title = "需要关注社交状态"
		analysis = fmt.Sprintf("您的社交回避及苦恼得分为%d分（总分28分），表明您有明显的社交困扰，可能影响日常生活。", totalScore)
		suggestions = []string{
			"系统学习社交技巧",
			"从小的社交目标开始练习",
			"考虑参加社交技能训练小组",
			"学习认知行为疗法技巧",
		}
	} else {
		level = "高度社交困扰"
		title = "建议寻求专业帮助"
		analysis = fmt.Sprintf("您的社交回避及苦恼得分为%d分（总分28分），表明您的社交困扰显著，建议寻求专业帮助。", totalScore)
		suggestions = []string{
			"寻求心理咨询师的专业帮助",
			"考虑认知行为疗法",
			"逐步暴露疗法可能有效",
			"建立支持系统，与理解的朋友交流",
		}
	}

	return map[string]interface{}{
		"test_name":   "社交回避及苦恼测试",
		"total_score": totalScore,
		"max_score":   28,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
