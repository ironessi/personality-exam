package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) Hostility(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// 敌意量表分级标准
	if totalScore >= 45 {
		level = "高敌意"
		title = "敌意水平较高"
		analysis = fmt.Sprintf("您的敌意水平得分为%d分（总分60分），表现出较高的敌意水平，可能经常感到愤怒、怨恨，对他人持怀疑态度。这种状态可能影响您的人际关系和心理健康。", totalScore)
		suggestions = []string{
			"学习情绪管理和愤怒控制技巧",
			"尝试认知行为疗法，改变负面思维模式",
			"练习放松技巧如深呼吸、冥想",
			"寻求专业心理咨询帮助",
			"培养宽容和理解他人的心态",
		}
	} else if totalScore >= 30 {
		level = "中等敌意"
		title = "敌意水平适中"
		analysis = fmt.Sprintf("您的敌意水平得分为%d分（总分60分），表明您有一定的敌意倾向，在特定情境下可能产生愤怒和怀疑情绪。需要关注情绪管理，避免影响人际关系。", totalScore)
		suggestions = []string{
			"改善沟通技巧，学习表达而非压抑情绪",
			"培养积极思维，减少猜疑",
			"参与团队合作活动，增进信任",
			"定期进行体育锻炼释放压力",
			"学习冲突解决方法",
		}
	} else {
		level = "低敌意"
		title = "敌意水平较低"
		analysis = fmt.Sprintf("您的敌意水平得分为%d分（总分60分），表明您的敌意水平较低，情绪控制能力较好，人际关系较为和谐。能够以平和的心态面对冲突和挫折。", totalScore)
		suggestions = []string{
			"继续保持良好的情绪管理习惯",
			"帮助他人改善人际关系和情绪控制",
			"参与社区服务活动，传播正能量",
			"定期自我反思，维持心理健康",
			"学习更高级的沟通和调解技巧",
		}
	}

	return map[string]interface{}{
		"test_name":   "愤怒-敌意测试",
		"total_score": totalScore,
		"max_score":   60,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
