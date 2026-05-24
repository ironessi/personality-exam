package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) UCLALoneliness(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// UCLA孤独量表分级标准
	if totalScore <= 20 {
		level = "低度孤独感"
		title = "孤独感轻微"
		analysis = fmt.Sprintf("您的孤独感得分为%d分（总分80分），表明您几乎没有孤独感，社交关系良好。", totalScore)
		suggestions = []string{
			"继续保持积极的社交活动和健康的生活方式",
		}
	} else if totalScore <= 34 {
		level = "轻度孤独感"
		title = "孤独感适中"
		analysis = fmt.Sprintf("您的孤独感得分为%d分（总分80分），表明您偶尔会感到孤独，但总体上适应良好。", totalScore)
		suggestions = []string{
			"主动维持社交联系",
			"培养兴趣爱好",
			"参与社区活动",
		}
	} else if totalScore <= 49 {
		level = "中度孤独感"
		title = "需要关注社交健康"
		analysis = fmt.Sprintf("您的孤独感得分为%d分（总分80分），表明您经常感到孤独，需要关注社交健康。", totalScore)
		suggestions = []string{
			"主动扩展社交圈",
			"寻求朋友支持",
			"参加兴趣小组",
		}
	} else if totalScore <= 59 {
		level = "中高度孤独感"
		title = "社交关系需要改善"
		analysis = fmt.Sprintf("您的孤独感得分为%d分（总分80分），表明您感到明显的孤独，社交关系需要改善。", totalScore)
		suggestions = []string{
			"参加社交活动",
			"考虑心理咨询",
			"建立支持系统",
		}
	} else {
		level = "高度孤独感"
		title = "需要专业支持"
		analysis = fmt.Sprintf("您的孤独感得分为%d分（总分80分），表明您感到强烈的孤独，需要专业支持。", totalScore)
		suggestions = []string{
			"寻求心理咨询师帮助",
			"建立支持系统",
			"参与团体治疗",
			"学习社交技能",
		}
	}

	return map[string]interface{}{
		"test_name":   "UCLA孤独量表",
		"total_score": totalScore,
		"max_score":   80,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
