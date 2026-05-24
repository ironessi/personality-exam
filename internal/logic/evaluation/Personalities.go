package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// PersonalitiesResult Personalities人格测试结果
type PersonalitiesResult struct {
	Extraversion      int      // 外向性
	Openness          int      // 开放性
	Conscientiousness int      // 尽责性
	Agreeableness     int      // 宜人性
	TotalScore        int      // 总分
	PersonalityType   string   // 人格类型
	Level             string   // 等级
	Title             string   // 标题
	Analysis          string   // 分析
	Suggestions       []string // 建议
}

// Personalities Personalities人格测试评分分析
func (e *Evaluation) Personalities(data []v1.AnswerRecord) (interface{}, error) {
	// 计算四个维度的得分
	extraversion := 0      // 外向性
	openness := 0          // 开放性
	conscientiousness := 0 // 尽责性
	agreeableness := 0     // 宜人性

	// 根据题目分配维度
	// 外向性（Extraversion）：题目1, 5, 9, 13, 17
	// 开放性（Openness）：题目2, 3, 7, 11, 15, 19
	// 尽责性（Conscientiousness）：题目6, 10, 14, 18
	// 宜人性（Agreeableness）：题目4, 8, 12, 16, 20

	// 简单的维度分配
	for _, record := range data {
		score := int(record.Score)
		switch record.QuestionID {
		case 1, 9, 17: // 外向性（正向计分）
			extraversion += score
		case 5, 13: // 外向性（反向计分）
			extraversion += 5 - score // 转换为正向计分
		case 2, 3, 7, 11, 15, 19: // 开放性
			openness += score
		case 6, 10, 14, 18: // 尽责性
			conscientiousness += score
		case 4, 8, 12, 16, 20: // 宜人性
			agreeableness += score
		}
	}

	totalScore := extraversion + openness + conscientiousness + agreeableness

	// 确定人格类型
	personalityType := determinePersonalitiesType(extraversion, openness, conscientiousness, agreeableness)

	// Personalities测试结果分析
	title := "Personalities人格测试结果"
	analysis := fmt.Sprintf("您的Personalities测试得分情况如下：\n外向性：%d分\n开放性：%d分\n尽责性：%d分\n宜人性：%d分\n总分：%d分\n\n您的人格类型：%s\n\n各维度解释：\n- 外向性：反映您的社交倾向，得分越高越外向，善于社交\n- 开放性：反映您对新事物的接受程度，得分越高越富有创造力和想象力\n- 尽责性：反映您的责任感和组织性，得分越高越可靠、有计划\n- 宜人性：反映您对他人的信任和友善程度，得分越高越容易与他人相处\n\n%s人格类型特点：\n- 这是一个基于您的测试得分确定的人格类型\n- 每种人格类型都有其独特的优势和特点\n- 没有绝对的好与坏，每种类型都有其适合的职业和生活方式", 
		extraversion, openness, conscientiousness, agreeableness, totalScore, personalityType, personalityType)

	// 提供建议
	suggestions := []string{
		"了解自己的人格特质有助于更好地理解自己和他人",
		"根据自己的人格特点选择适合的职业和生活方式",
		"在工作和生活中发挥自己的人格优势",
		"如果某些维度得分过低，建议尝试适当调整自己的行为方式",
		"定期进行人格测试，了解自己的人格变化",
	}

	return map[string]interface{}{
		"test_name":          "Personalities人格测试",
		"extraversion":       extraversion,
		"openness":           openness,
		"conscientiousness":  conscientiousness,
		"agreeableness":      agreeableness,
		"total_score":        totalScore,
		"personality_type":   personalityType,
		"level":              "完成测试",
		"title":              title,
		"analysis":           analysis,
		"suggestions":        suggestions,
	}, nil
}

// determinePersonalitiesType 确定Personalities人格类型
func determinePersonalitiesType(extraversion, openness, conscientiousness, agreeableness int) string {
	// 找出最高得分的维度
	dimensions := map[string]int{
		"外向性":  extraversion,
		"开放性":  openness,
		"尽责性":  conscientiousness,
		"宜人性":  agreeableness,
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
	case "开放性":
		return "开放型人格"
	case "尽责性":
		return "尽责型人格"
	case "宜人性":
		return "宜人型人格"
	default:
		return "综合型人格"
	}
}