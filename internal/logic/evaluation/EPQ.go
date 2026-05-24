package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// EPQ EPQ艾森克人格测试评分分析
func (e *Evaluation) EPQ(data []v1.AnswerRecord) (interface{}, error) {
	// 计算EPQ四个维度的得分
	// 这里简单地将题目分为四个维度，实际应用中需要根据具体的EPQ测试题目分配来计算
	neuroticism := 0 // 神经质(N)
	extraversion := 0  // 外向性(E)
	psychoticism := 0  // 精神质(P)
	lying := 0         // 说谎(L)

	// 简单的维度分配（实际应用中需要根据具体题目调整）
	for i, record := range data {
		score := int(record.Score)
		switch i % 4 {
		case 0:
			neuroticism += score
		case 1:
			extraversion += score
		case 2:
			psychoticism += score
		case 3:
			lying += score
		}
	}

	// EPQ人格测试结果分析
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	title = "EPQ艾森克人格测试结果"
	analysis = fmt.Sprintf("您的EPQ测试得分情况如下：\n神经质(N)：%d分\n外向性(E)：%d分\n精神质(P)：%d分\n说谎(L)：%d分\n\n各维度解释：\n- 神经质(N)：反映情绪稳定性，得分越高情绪越不稳定，容易焦虑、抑郁等。\n- 外向性(E)：反映性格内外向，得分越高越外向，善于社交。\n- 精神质(P)：反映心理变态倾向，得分越高越孤独、冷酷、敌视等。\n- 说谎(L)：反映测试的可靠性，得分过高可能表示测试结果不可靠。", neuroticism, extraversion, psychoticism, lying)
	level = "完成测试"

	// 提供建议
	suggestions = []string{
		"了解自己的人格特质有助于更好地理解自己和他人",
		"根据自己的人格特点选择适合的职业和生活方式",
		"如果某些维度得分过高，建议寻求专业心理咨询师的帮助",
		"定期进行人格测试，了解自己的人格变化",
	}

	return map[string]interface{}{
		"test_name":    "EPQ艾森克人格测试-国际标准版",
		"neuroticism":  neuroticism,
		"extraversion": extraversion,
		"psychoticism": psychoticism,
		"lying":        lying,
		"level":        level,
		"title":        title,
		"analysis":     analysis,
		"suggestions":  suggestions,
	}, nil
}