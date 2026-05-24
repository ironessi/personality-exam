package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) PROMISAnxiety(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// PROMIS焦虑量表分级标准
	if totalScore <= 35 {
		level = "轻度焦虑"
		title = "焦虑症状轻微"
		analysis = fmt.Sprintf("您的焦虑症状得分为%d分（总分105分），表明您有轻微的焦虑症状，属于正常范围内的情绪波动。", totalScore)
		suggestions = []string{
			"保持规律作息和适度运动",
			"学习简单的放松技巧",
			"维持健康的社交活动",
		}
	} else if totalScore <= 60 {
		level = "中度焦虑"
		title = "焦虑症状明显"
		analysis = fmt.Sprintf("您的焦虑症状得分为%d分（总分105分），表明您有明显的焦虑症状，可能影响日常生活。", totalScore)
		suggestions = []string{
			"学习压力管理技巧",
			"练习深呼吸和冥想",
			"适当减少压力源",
			"考虑心理咨询",
		}
	} else if totalScore <= 85 {
		level = "明显焦虑"
		title = "焦虑症状显著"
		analysis = fmt.Sprintf("您的焦虑症状得分为%d分（总分105分），表明您有显著的焦虑症状，建议积极应对。", totalScore)
		suggestions = []string{
			"寻求专业心理咨询",
			"学习认知行为疗法技巧",
			"建立支持系统",
			"必要时就医评估",
		}
	} else {
		level = "严重焦虑"
		title = "焦虑症状严重"
		analysis = fmt.Sprintf("您的焦虑症状得分为%d分（总分105分），表明您有严重的焦虑症状，强烈建议寻求专业帮助。", totalScore)
		suggestions = []string{
			"立即寻求心理医生帮助",
			"考虑药物治疗配合心理治疗",
			"建立紧急支持网络",
			"学习危机应对策略",
		}
	}

	return map[string]interface{}{
		"test_name":   "焦虑性测试",
		"total_score": totalScore,
		"max_score":   105,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
