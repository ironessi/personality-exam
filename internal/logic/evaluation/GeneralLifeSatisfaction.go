package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) GeneralLifeSatisfaction(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	userAnswers := make([]int, len(data))
	for i, record := range data {
		score := int(record.Score)
		totalScore += score
		userAnswers[i] = score
	}

	// 调用提供反馈的函数（可以改为返回结果）
	maxScore := 30 // 假设10题，每题1-3分，根据实际情况调整
	result := provideFeedbackResult(totalScore, maxScore, userAnswers)

	return result, nil
}

// 重构 provideFeedback 为返回结果而不是直接打印
func provideFeedbackResult(totalScore, maxScore int, userAnswers []int) map[string]interface{} {
	percentage := float64(totalScore) / float64(maxScore) * 100

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	if percentage >= 80 {
		level = "非常满意"
		title = "生活满意度很高"
		analysis = "🎉 恭喜！您的生活满意度很高！"
		suggestions = []string{
			"继续保持积极的生活态度和健康的生活方式",
			"与亲友分享您的快乐，传递正能量",
			"设定新的目标，持续成长",
		}
	} else if percentage >= 60 {
		level = "比较满意"
		title = "生活满意度中等"
		analysis = "👍 您的生活满意度处于中等水平"
		suggestions = []string{
			"关注生活中的积极方面，培养感恩的心态",
			"加强社交联系，与朋友家人多交流",
			"寻找新的兴趣爱好，丰富生活体验",
		}
	} else if percentage >= 40 {
		level = "一般"
		title = "生活满意度有待提升"
		analysis = "🤔 您的生活满意度有待提升"
		suggestions = []string{
			"尝试改变日常routine，寻找新的生活乐趣",
			"考虑与专业人士交流，获得心理支持",
			"设定小目标，逐步改善生活质量",
			"多参与户外活动，接触大自然",
		}
	} else {
		level = "不满意"
		title = "面临生活挑战"
		analysis = "💡 您可能正面临一些挑战"
		suggestions = []string{
			"寻求专业心理咨询师的帮助",
			"与信任的亲友倾诉您的感受",
			"关注身心健康，保证充足的休息",
			"记住：困难是暂时的，您并不孤单",
			"如有需要，可拨打心理援助热线",
		}
	}

	// 添加个性化提醒
	var personalizedReminders []string
	if len(userAnswers) > 6 && userAnswers[6] == 3 { // 经常感到孤独
		personalizedReminders = append(personalizedReminders, "您表示经常感到孤独，建议主动参加社交活动或兴趣小组")
	}
	if len(userAnswers) > 7 && userAnswers[7] == 3 { // 经常感到生活无目的
		personalizedReminders = append(personalizedReminders, "您感到生活缺乏目的，可以尝试志愿服务或学习新技能")
	}
	if len(userAnswers) > 4 && userAnswers[4] == 3 { // 担心无法满足他人期望
		personalizedReminders = append(personalizedReminders, "您担心无法满足他人期望，请记住：做最好的自己就足够了")
	}
	if len(userAnswers) > 9 && userAnswers[9] == 3 { // 发现许多不幸的事情
		personalizedReminders = append(personalizedReminders, "您发现生活中许多不幸，建议关注微小但美好的事物")
	}

	return map[string]interface{}{
		"test_name":              "生活满意度测试",
		"total_score":            totalScore,
		"max_score":              maxScore,
		"percentage":             fmt.Sprintf("%.1f%%", percentage),
		"level":                  level,
		"title":                  title,
		"analysis":               analysis,
		"suggestions":            suggestions,
		"personalized_reminders": personalizedReminders,
	}
}
