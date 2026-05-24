package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// DISC DISC性格测试评分分析
func (e *Evaluation) DISC(data []v1.AnswerRecord) (interface{}, error) {
	// 计算DISC四个维度的得分
	// 这里简单地将题目分为四个维度，实际应用中需要根据具体的DISC测试题目分配来计算
	dScore := 0
	iScore := 0
	sScore := 0
	cScore := 0

	// 简单的维度分配（实际应用中需要根据具体题目调整）
	for i, record := range data {
		score := int(record.Score)
		switch i % 4 {
		case 0:
			dScore += score
		case 1:
			iScore += score
		case 2:
			sScore += score
		case 3:
			cScore += score
		}
	}

	// 确定主导性格类型
	maxScore := dScore
	dominantType := "支配型(D)"
	if iScore > maxScore {
		maxScore = iScore
		dominantType = "影响型(I)"
	}
	if sScore > maxScore {
		maxScore = sScore
		dominantType = "稳健型(S)"
	}
	if cScore > maxScore {
		dominantType = "谨慎型(C)"
	}

	// DISC性格测试结果分析
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	title = dominantType + "性格特质"
	analysis = fmt.Sprintf("您的主导性格类型是%s，具体得分情况如下：\n支配型(D)：%d分\n影响型(I)：%d分\n稳健型(S)：%d分\n谨慎型(C)：%d分", dominantType, dScore, iScore, sScore, cScore)
	level = "完成测试"

	// 根据主导性格类型提供建议
	switch dominantType {
	case "支配型(D)":
		suggestions = []string{
			"注意倾听他人的意见，避免过于独断专行",
			"学会控制自己的脾气，避免冲动决策",
			"关注团队合作，培养协作能力",
			"学会赞美和肯定他人",
		}
	case "影响型(I)":
		suggestions = []string{
			"注意细节，避免粗心大意",
			"学会专注，避免同时处理过多事情",
			"培养耐心，学会倾听他人",
			"制定计划，提高执行力",
		}
	case "稳健型(S)":
		suggestions = []string{
			"尝试接受变化，培养适应能力",
			"学会表达自己的意见和需求",
			"培养决策能力，避免过于依赖他人",
			"尝试挑战自己，走出舒适区",
		}
	case "谨慎型(C)":
		suggestions = []string{
			"学会灵活变通，避免过于刻板",
			"培养团队合作精神，学会信任他人",
			"学会放松，避免过度追求完美",
			"提高沟通能力，学会表达自己的想法",
		}
	}

	return map[string]interface{}{
		"test_name":     "DISC性格测试",
		"d_score":       dScore,
		"i_score":       iScore,
		"s_score":       sScore,
		"c_score":       cScore,
		"dominant_type": dominantType,
		"level":         level,
		"title":         title,
		"analysis":      analysis,
		"suggestions":   suggestions,
	}, nil
}
