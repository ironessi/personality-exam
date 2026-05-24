package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) SelfEfficacy(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// 自我效能感量表分级标准
	if totalScore >= 40 {
		level = "很高"
		title = "自我效能感很强"
		analysis = fmt.Sprintf("您的自我效能感得分为%d分（总分50分），表明您对自己解决问题的能力充满信心，面对挑战时能够保持积极态度，相信自己能够克服困难并达成目标。", totalScore)
		suggestions = []string{
			"继续保持自信和积极的心态",
			"挑战更高难度的目标，持续成长",
			"成为他人的榜样，帮助建立自信",
			"反思成功经验，强化自我效能",
			"在团队中发挥领导作用",
		}
	} else if totalScore >= 30 {
		level = "中等"
		title = "自我效能感适中"
		analysis = fmt.Sprintf("您的自我效能感得分为%d分（总分50分），表明您有一定的自信心，但在面对较大挑战时可能产生怀疑。通过积累成功经验可以进一步提升自我效能感。", totalScore)
		suggestions = []string{
			"通过实现小目标积累成功经验",
			"学习问题解决和决策技巧",
			"观察和学习他人的成功经验",
			"设定具体可行的目标计划",
			"寻求反馈和支持，增强信心",
		}
	} else {
		level = "较低"
		title = "自我效能感需提升"
		analysis = fmt.Sprintf("您的自我效能感得分为%d分（总分50分），表明您可能需要提升自信心和问题解决能力，面对挑战时容易感到无助或怀疑自己的能力。", totalScore)
		suggestions = []string{
			"设定可实现的小目标，逐步建立信心",
			"记录和庆祝每一个小成功",
			"寻求他人的鼓励和支持",
			"学习压力管理和放松技巧",
			"参加自信心训练或技能培训",
		}
	}

	return map[string]interface{}{
		"test_name":   "自我效能感测试",
		"total_score": totalScore,
		"max_score":   50,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}
