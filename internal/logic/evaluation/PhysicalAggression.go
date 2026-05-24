package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) PhysicalAggression(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// 身体攻击量表分级标准
	if totalScore >= 15 {
		level = "高风险"
		title = "攻击倾向较高"
		analysis = fmt.Sprintf("您的攻击倾向得分为%d分（总分19分），表现出较高的攻击行为倾向，可能容易冲动，缺乏自我控制，在冲突情境中倾向于使用攻击性行为。", totalScore)
		suggestions = []string{
			"立即寻求专业心理咨询或治疗",
			"学习愤怒管理和冲动控制技巧",
			"避免接触暴力性媒体内容",
			"培养替代性的压力释放方式如运动",
			"建立支持系统，与理解的人交流",
		}
	} else if totalScore >= 10 {
		level = "中等风险"
		title = "攻击倾向适中"
		analysis = fmt.Sprintf("您的攻击倾向得分为%d分（总分19分），表明您有一定的攻击倾向，在压力或冲突情境中可能表现出攻击性行为。需要关注情绪调节和冲动控制。", totalScore)
		suggestions = []string{
			"学习情绪识别和表达技巧",
			"练习放松技巧如深呼吸、冥想",
			"培养解决问题的替代方法",
			"参与体育活动释放能量",
			"学习有效的沟通和冲突解决技巧",
		}
	} else {
		level = "低风险"
		title = "攻击倾向较低"
		analysis = fmt.Sprintf("您的攻击倾向得分为%d分（总分19分），表明您的攻击倾向较低，具有良好的自我控制能力，能够以建设性的方式处理冲突和挫折。", totalScore)
		suggestions = []string{
			"继续保持良好的情绪管理",
			"帮助他人学习冲突解决技巧",
			"参与和平建设或调解活动",
			"定期进行心理健康自我评估",
			"培养同理心和理解他人的能力",
		}
	}

	return map[string]interface{}{
		"test_name":   "愤怒-身体攻击测试",
		"total_score": totalScore,
		"max_score":   19,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
