package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// RIASECResult 霍兰德职业兴趣测试(RIASEC)结果
type RIASECResult struct {
	Realistic       int      // 现实型(R)
	Investigative   int      // 研究型(I)
	Artistic        int      // 艺术型(A)
	Social          int      // 社会型(S)
	Enterprising    int      // 企业型(E)
	Conventional    int      // 常规型(C)
	TotalScore      int      // 总分
	DominantType    string   // 主导兴趣类型
	SecondaryType   string   // 次要兴趣类型
	Level           string   // 等级
	Title           string   // 标题
	Analysis        string   // 分析
	Suggestions     []string // 建议
}

// RIASEC 霍兰德职业兴趣测试(RIASEC)评分分析
func (e *Evaluation) RIASEC(data []v1.AnswerRecord) (interface{}, error) {
	// 计算六种类型的得分
	realistic := 0       // 现实型(R)
	investigative := 0   // 研究型(I)
	artistic := 0        // 艺术型(A)
	social := 0          // 社会型(S)
	enterprising := 0    // 企业型(E)
	conventional := 0    // 常规型(C)

	// 根据题目分配维度
	// 现实型(R)：题目1, 7, 13, 19, 25, 31, 37, 43, 49, 55, 61, 67, 73, 79, 85, 91, 97, 103, 109
	// 研究型(I)：题目2, 8, 14, 20, 26, 32, 38, 44, 50, 56, 62, 68, 74, 80, 86, 92, 98, 104, 110
	// 艺术型(A)：题目3, 9, 15, 21, 27, 33, 39, 45, 51, 57, 63, 69, 75, 81, 87, 93, 99, 105
	// 社会型(S)：题目4, 10, 16, 22, 28, 34, 40, 46, 52, 58, 64, 70, 76, 82, 88, 94, 100, 106
	// 企业型(E)：题目5, 11, 17, 23, 29, 35, 41, 47, 53, 59, 65, 71, 77, 83, 89, 95, 101, 107
	// 常规型(C)：题目6, 12, 18, 24, 30, 36, 42, 48, 54, 60, 66, 72, 78, 84, 90, 96, 102, 108

	// 简单的维度分配
	for _, record := range data {
		score := int(record.Score)
		questionID := record.QuestionID

		// 现实型(R)
		if questionID%6 == 1 {
			realistic += score
		} else if questionID%6 == 2 {
			// 研究型(I)
			investigative += score
		} else if questionID%6 == 3 {
			// 艺术型(A)
			artistic += score
		} else if questionID%6 == 4 {
			// 社会型(S)
			social += score
		} else if questionID%6 == 5 {
			// 企业型(E)
			enterprising += score
		} else if questionID%6 == 0 {
			// 常规型(C)
			conventional += score
		}
	}

	totalScore := realistic + investigative + artistic + social + enterprising + conventional

	// 确定主导兴趣类型和次要兴趣类型
	dominantType, secondaryType := determineRIASECTypes(realistic, investigative, artistic, social, enterprising, conventional)

	// RIASEC测试结果分析
	title := "霍兰德职业兴趣测试(RIASEC)结果"
	analysis := fmt.Sprintf("您的RIASEC测试得分情况如下：\n现实型(R)：%d分\n研究型(I)：%d分\n艺术型(A)：%d分\n社会型(S)：%d分\n企业型(E)：%d分\n常规型(C)：%d分\n总分：%d分\n\n主导兴趣类型：%s\n次要兴趣类型：%s\n\n各类型解释：\n- 现实型(R)：喜欢使用工具或机器，从事实际操作类工作\n- 研究型(I)：喜欢思考和解决复杂问题，从事科学或研究类工作\n- 艺术型(A)：喜欢创作和表达，从事艺术或创意类工作\n- 社会型(S)：喜欢帮助他人，从事教育或服务类工作\n- 企业型(E)：喜欢领导和管理，从事商务或管理类工作\n- 常规型(C)：喜欢按规则和程序工作，从事行政或事务类工作\n\n%s和%s兴趣类型特点：\n- %s类型的人通常具有...的特点\n- %s类型的人通常具有...的特点\n\n适合的职业方向：\n- 适合从事与%s和%s相关的职业\n- 具体职业包括：...等", 
		realistic, investigative, artistic, social, enterprising, conventional, totalScore, dominantType, secondaryType, dominantType, secondaryType, dominantType, secondaryType, dominantType, secondaryType)

	// 提供建议
	suggestions := []string{
		"了解自己的职业兴趣类型有助于选择适合的职业方向",
		"根据兴趣类型探索相关职业，了解职业的具体内容和要求",
		"考虑将兴趣与能力相结合，选择最适合自己的职业",
		"定期进行职业兴趣测试，追踪自己的兴趣变化",
		"参加相关的职业培训或实践活动，提升自己的职业能力",
	}

	return map[string]interface{}{
		"test_name":      "霍兰德职业兴趣测试（RIASEC）",
		"realistic":      realistic,
		"investigative":  investigative,
		"artistic":       artistic,
		"social":         social,
		"enterprising":   enterprising,
		"conventional":   conventional,
		"total_score":    totalScore,
		"dominant_type":  dominantType,
		"secondary_type": secondaryType,
		"level":          "完成测试",
		"title":          title,
		"analysis":       analysis,
		"suggestions":    suggestions,
	}, nil
}

// determineRIASECTypes 确定RIASEC主导兴趣类型和次要兴趣类型
func determineRIASECTypes(realistic, investigative, artistic, social, enterprising, conventional int) (string, string) {
	// 定义类型映射
	types := map[string]int{
		"现实型(R)": realistic,
		"研究型(I)": investigative,
		"艺术型(A)": artistic,
		"社会型(S)": social,
		"企业型(E)": enterprising,
		"常规型(C)": conventional,
	}

	// 找出最高得分和次高得分的类型
	maxScore1 := 0
	maxType1 := ""
	maxScore2 := 0
	maxType2 := ""

	for typeName, score := range types {
		if score > maxScore1 {
			maxScore2 = maxScore1
			maxType2 = maxType1
			maxScore1 = score
			maxType1 = typeName
		} else if score > maxScore2 {
			maxScore2 = score
			maxType2 = typeName
		}
	}

	return maxType1, maxType2
}