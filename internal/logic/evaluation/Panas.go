package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) Panas(data []v1.AnswerRecord) (interface{}, error) {
	positiveScore := 0
	negativeScore := 0
	totalScore := 0

	// 积极情感项目：第6,7,8,9,10,14,15,16,17题（对应索引5-9,13-16）
	positiveItems := []int{5, 6, 7, 8, 9, 13, 14, 15, 16}
	// 消极情感项目：第1,2,3,4,5,11,12,13,18题（对应索引0-4,10-12,17）
	negativeItems := []int{0, 1, 2, 3, 4, 10, 11, 12, 17}

	// 计算总分和各维度分
	for i, record := range data {
		score := int(record.Score)
		totalScore += score

		// 检查是否是积极情感项目
		for _, idx := range positiveItems {
			if i == idx {
				positiveScore += score
				break
			}
		}

		// 检查是否是消极情感项目
		for _, idx := range negativeItems {
			if i == idx {
				negativeScore += score
				break
			}
		}
	}

	// 计算情感平衡分数
	affectBalance := positiveScore - negativeScore
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// 判断等级
	if affectBalance > 10 {
		level = "高度积极"
		title = "情感状态非常良好"
		analysis = fmt.Sprintf("您的积极情感(%d分)明显多于消极情感(%d分)，情感平衡分数为%d分。情绪状态非常健康，您可能经常体验到快乐、满足和活力。", positiveScore, negativeScore, affectBalance)
		suggestions = []string{
			"继续保持积极的生活方式和心态",
			"多参与让您感到快乐和满足的活动",
			"与积极乐观的人保持交往",
			"记录和分享积极体验",
			"帮助他人培养积极情感",
		}
	} else if affectBalance > 0 {
		level = "适度积极"
		title = "情感状态良好"
		analysis = fmt.Sprintf("您的积极情感(%d分)多于消极情感(%d分)，情感平衡分数为%d分。情绪状态总体良好，能够较好地处理日常压力。", positiveScore, negativeScore, affectBalance)
		suggestions = []string{
			"继续培养积极情绪体验",
			"学习压力管理技巧",
			"保持规律的运动和休息",
			"培养感恩的心态",
			"建立支持性的社交网络",
		}
	} else if affectBalance == 0 {
		level = "平衡"
		title = "情感状态平衡"
		analysis = fmt.Sprintf("您的积极情感(%d分)和消极情感(%d分)相对平衡，情感平衡分数为%d分。可能需要更多关注情绪调节和积极体验的培养。", positiveScore, negativeScore, affectBalance)
		suggestions = []string{
			"增加积极活动的参与",
			"练习正念和情绪觉察",
			"学习情绪调节策略",
			"设定有意义的目标",
			"寻求社交支持和连接",
		}
	} else {
		level = "需要关注"
		title = "情感状态需关注"
		analysis = fmt.Sprintf("您的消极情感(%d分)较多，积极情感(%d分)相对较少，情感平衡分数为%d分。可能经常体验到压力、焦虑或低落情绪，建议关注情绪健康。", negativeScore, positiveScore, affectBalance)
		suggestions = []string{
			"尝试正念冥想和放松技巧",
			"寻求朋友和家人的支持",
			"考虑专业心理咨询",
			"增加体育锻炼和户外活动",
			"培养兴趣爱好和愉悦活动",
		}
	}

	return map[string]interface{}{
		"test_name":      "积极情感消极情感量表(PANAS)",
		"total_score":    totalScore,
		"positive_score": positiveScore,
		"negative_score": negativeScore,
		"affect_balance": affectBalance,
		"max_score":      90,
		"level":          level,
		"title":          title,
		"analysis":       analysis,
		"suggestions":    suggestions,
	}, nil
}
