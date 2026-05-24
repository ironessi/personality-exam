package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// PhilPersonalityResult 菲尔人格十项测试题结果
type PhilPersonalityResult struct {
	Extraversion      int      // 外向性
	Agreeableness     int      // 宜人性
	Conscientiousness int      // 尽责性
	EmotionalStability int     // 情绪稳定性
	Openness          int      // 开放性
	TotalScore        int      // 总分
	Level             string   // 等级
	Title             string   // 标题
	Analysis          string   // 分析
	Suggestions       []string // 建议
	PersonalityType   string   // 人格类型
}

// PhilPersonality 菲尔人格十项测试题评分分析
func (e *Evaluation) PhilPersonality(data []v1.AnswerRecord) (interface{}, error) {
	// 计算五个维度的得分
	extraversion := 0      // 外向性
	agreeableness := 0     // 宜人性
	conscientiousness := 0 // 尽责性
	emotionalStability := 0 // 情绪稳定性
	openness := 0          // 开放性

	// 根据题目分配维度
	// 题目1: 我喜欢参加社交活动，认识新朋友 - 外向性
	// 题目2: 我倾向于信任他人，认为大多数人是善良的 - 宜人性
	// 题目3: 我做事认真负责，总是尽力完成任务 - 尽责性
	// 题目4: 我情绪稳定，很少大起大落 - 情绪稳定性
	// 题目5: 我喜欢尝试新事物，对未知充满好奇 - 开放性
	// 题目6: 我喜欢独立思考，不随波逐流 - 开放性
	// 题目7: 我对他人的感受很敏感，善于理解他人 - 开放性
	// 题目8: 我喜欢有计划地安排生活，不喜欢混乱 - 尽责性
	// 题目9: 我面对压力时能保持冷静，有效应对 - 情绪稳定性
	// 题目10: 我对艺术、音乐和文学有浓厚兴趣 - 开放性

	// 简单的维度分配
	for _, record := range data {
		score := int(record.Score)
		switch record.QuestionID {
		case 1: // 外向性
			extraversion += score
		case 2: // 宜人性
			agreeableness += score
		case 3, 8: // 尽责性
			conscientiousness += score
		case 4, 9: // 情绪稳定性
			emotionalStability += score
		case 5, 6, 7, 10: // 开放性
			openness += score
		}
	}

	totalScore := extraversion + agreeableness + conscientiousness + emotionalStability + openness

	// 确定人格类型
	personalityType := determinePersonalityType(extraversion, agreeableness, conscientiousness, emotionalStability, openness)

	// 菲尔人格测试结果分析
	title := "菲尔人格十项测试题结果"
	analysis := fmt.Sprintf("您的菲尔人格测试得分情况如下：\n外向性：%d分\n宜人性：%d分\n尽责性：%d分\n情绪稳定性：%d分\n开放性：%d分\n总分：%d分\n\n您的人格类型：%s\n\n各维度解释：\n- 外向性：反映您的社交倾向，得分越高越外向，善于社交。\n- 宜人性：反映您对他人的信任和友善程度，得分越高越容易与他人相处。\n- 尽责性：反映您的责任感和组织性，得分越高越可靠、有计划。\n- 情绪稳定性：反映您的情绪控制能力，得分越高情绪越稳定。\n- 开放性：反映您对新事物的接受程度，得分越高越富有创造力和想象力。", 
		extraversion, agreeableness, conscientiousness, emotionalStability, openness, totalScore, personalityType)

	// 提供建议
	suggestions := []string{
		"了解自己的人格特质有助于更好地理解自己和他人",
		"根据自己的人格特点选择适合的职业和生活方式",
		"在工作和生活中发挥自己的人格优势",
		"如果某些维度得分过低，建议尝试适当调整自己的行为方式",
		"定期进行人格测试，了解自己的人格变化",
	}

	return map[string]interface{}{
		"test_name":          "菲尔人格十项测试题",
		"extraversion":       extraversion,
		"agreeableness":      agreeableness,
		"conscientiousness":  conscientiousness,
		"emotional_stability": emotionalStability,
		"openness":           openness,
		"total_score":        totalScore,
		"level":              "完成测试",
		"title":              title,
		"analysis":           analysis,
		"suggestions":        suggestions,
		"personality_type":   personalityType,
	}, nil
}

// determinePersonalityType 确定人格类型
func determinePersonalityType(extraversion, agreeableness, conscientiousness, emotionalStability, openness int) string {
	// 找出最高得分的维度
	dimensions := map[string]int{
		"外向性":  extraversion,
		"宜人性":  agreeableness,
		"尽责性":  conscientiousness,
		"情绪稳定性": emotionalStability,
		"开放性":  openness,
	}

	maxDimension := ""
	maxScore := 0
	for dim, score := range dimensions {
		if score > maxScore {
			maxScore = score
			maxDimension = dim
		}
	}

	// 根据最高维度确定人格类型
	switch maxDimension {
	case "外向性":
		return "外向型人格"
	case "宜人性":
		return "宜人型人格"
	case "尽责性":
		return "尽责型人格"
	case "情绪稳定性":
		return "稳定型人格"
	case "开放性":
		return "开放型人格"
	default:
		return "综合型人格"
	}
}