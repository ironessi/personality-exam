package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// MBTIResult 16型人格测试(MBTI)结果
type MBTIResult struct {
	Extraversion    int      // 外倾(E)vs内倾(I)
	Sensing         int      // 感觉(S)vs直觉(N)
	Thinking        int      // 思考(T)vs情感(F)
	Judging         int      // 判断(J)vs感知(P)
	TotalScore      int      // 总分
	PersonalityType string   // 人格类型
	Level           string   // 等级
	Title           string   // 标题
	Analysis        string   // 分析
	Suggestions     []string // 建议
}

// MBTI 16型人格测试(MBTI)评分分析
func (e *Evaluation) MBTI(data []v1.AnswerRecord) (interface{}, error) {
	// 计算四个维度的得分
	extraversion := 0 // 外倾(E)vs内倾(I)
	sensing := 0      // 感觉(S)vs直觉(N)
	thinking := 0     // 思考(T)vs情感(F)
	judging := 0      // 判断(J)vs感知(P)

	// 根据题目分配维度
	// 题目1: 我更倾向于通过社交活动获取能量 - 外倾(E)vs内倾(I)
	// 题目2: 我更关注事实和细节，而非抽象概念 - 感觉(S)vs直觉(N)
	// 题目3: 我做决策时更依赖逻辑分析，而非情感因素 - 思考(T)vs情感(F)
	// 题目4: 我喜欢计划和组织，倾向于有秩序的生活 - 判断(J)vs感知(P)
	// 题目5: 我喜欢与他人合作，而非独立工作 - 外倾(E)vs内倾(I)
	// 题目6: 我更相信经验，而非直觉 - 感觉(S)vs直觉(N)
	// 题目7: 我更注重公平，而非和谐 - 思考(T)vs情感(F)
	// 题目8: 我喜欢提前安排好事情，而非即兴发挥 - 判断(J)vs感知(P)
	// 题目9: 我喜欢在大群体中社交 - 外倾(E)vs内倾(I)
	// 题目10: 我更关注当前，而非未来 - 感觉(S)vs直觉(N)

	// 题目 ID → 维度：用「相对题号」而不是硬编码的全局 ID。
	//
	// 这里原本写的是 switch record.QuestionID { case 1, 5, 9: ... }，依赖题目
	// 自增 ID 恰好从 1 开始。用 public.sql 建库再走 /api/initialize 时，MBTI 的
	// 题目 ID 是 376..385，四个维度全部累加为 0，任何作答都会判成 INFP。
	// 现在由 Evaluate 入口的 normalizeQuestionIDs 统一折算成 1..N（见
	// evaluation.go），所以这里继续按题号 1..10 写即可。
	// 前提是题目 ID 连续且按题目顺序作答；缺题会让后续题号前移。

	// 简单的维度分配
	for _, record := range data {
		score := int(record.Score)
		switch record.QuestionID {
		case 1, 5, 9: // 外倾(E)vs内倾(I)
			extraversion += score
		case 2, 6, 10: // 感觉(S)vs直觉(N)
			sensing += score
		case 3, 7: // 思考(T)vs情感(F)
			thinking += score
		case 4, 8: // 判断(J)vs感知(P)
			judging += score
		}
	}

	totalScore := extraversion + sensing + thinking + judging

	// 确定人格类型
	personalityType := determineMBTIType(extraversion, sensing, thinking, judging)

	// MBTI测试结果分析
	title := "16型人格测试(MBTI)结果"
	analysis := fmt.Sprintf("您的MBTI测试得分情况如下：\n外倾(E)/内倾(I)：%d分\n感觉(S)/直觉(N)：%d分\n思考(T)/情感(F)：%d分\n判断(J)/感知(P)：%d分\n总分：%d分\n\n您的人格类型：%s\n\n各维度解释：\n- 外倾(E)vs内倾(I)：反映您获取能量的方式，外倾者通过社交活动获取能量，内倾者通过独处获取能量\n- 感觉(S)vs直觉(N)：反映您收集信息的方式，感觉者关注事实和细节，直觉者关注抽象概念和可能性\n- 思考(T)vs情感(F)：反映您做决策的方式，思考者依赖逻辑分析，情感者依赖情感因素\n- 判断(J)vs感知(P)：反映您组织生活的方式，判断者喜欢计划和组织，感知者喜欢灵活和即兴\n\n%s人格类型特点：\n- 这是一个基于您的测试得分确定的人格类型\n- 每种MBTI类型都有其独特的优势和特点\n- 没有绝对的好与坏，每种类型都有其适合的职业和生活方式",
		extraversion, sensing, thinking, judging, totalScore, personalityType, personalityType)

	// 提供建议
	suggestions := []string{
		"了解自己的MBTI类型有助于更好地理解自己和他人",
		"根据自己的人格特点选择适合的职业和生活方式",
		"在工作和生活中发挥自己的人格优势",
		"学习与不同人格类型的人有效沟通",
		"定期进行MBTI测试，了解自己的人格变化",
	}

	return map[string]interface{}{
		"test_name":        "16型人格测试（MBTI）",
		"extraversion":     extraversion,
		"sensing":          sensing,
		"thinking":         thinking,
		"judging":          judging,
		"total_score":      totalScore,
		"personality_type": personalityType,
		"level":            "完成测试",
		"title":            title,
		"analysis":         analysis,
		"suggestions":      suggestions,
	}, nil
}

// determineMBTIType 确定MBTI人格类型
func determineMBTIType(extraversion, sensing, thinking, judging int) string {
	// 确定每个维度的偏好
	// 外倾(E)vs内倾(I)：得分越高越外倾
	// 感觉(S)vs直觉(N)：得分越高越感觉
	// 思考(T)vs情感(F)：得分越高越思考
	// 判断(J)vs感知(P)：得分越高越判断

	var eorI, sorN, torF, jorP string

	// 外倾(E)vs内倾(I)：阈值为12（3个问题，每个问题7分，总分21，阈值12）
	if extraversion >= 12 {
		eorI = "E"
	} else {
		eorI = "I"
	}

	// 感觉(S)vs直觉(N)：阈值为12（3个问题，每个问题7分，总分21，阈值12）
	if sensing >= 12 {
		sorN = "S"
	} else {
		sorN = "N"
	}

	// 思考(T)vs情感(F)：阈值为8（2个问题，每个问题7分，总分14，阈值8）
	if thinking >= 8 {
		torF = "T"
	} else {
		torF = "F"
	}

	// 判断(J)vs感知(P)：阈值为8（2个问题，每个问题7分，总分14，阈值8）
	if judging >= 8 {
		jorP = "J"
	} else {
		jorP = "P"
	}

	// 组合四个维度，形成MBTI人格类型
	return eorI + sorN + torF + jorP
}
