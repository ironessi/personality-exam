package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
)

func (e *Evaluation) PerceivedStress(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0
	userAnswers := make([]int, len(data))
	for i, record := range data {
		score := int(record.Score)
		totalScore += score
		userAnswers[i] = score
	}

	// 生成结果
	result := generateStressResult(totalScore, userAnswers)
	return result, nil
}

func generateStressResult(totalScore int, userAnswers []int) map[string]interface{} {
	level := ""
	title := ""
	analysis := ""
	var suggestions []string
	var highScoreItems []string

	// 压力程度分析
	if totalScore <= 33 {
		level = "低度压力"
		title = "压力水平较低"
		analysis = "您的压力水平较低，能够较好地应对生活挑战。"
	} else if totalScore <= 66 {
		level = "中度压力"
		title = "压力水平适中"
		analysis = "您感受到中等程度的压力，需要适当关注和管理。"
	} else {
		level = "高度压力"
		title = "压力水平较高"
		analysis = "您的压力水平较高，建议积极寻求改善方法。"
	}

	// 识别高分项目（主要压力源）
	highScoreItemDescs := []struct {
		index int
		desc  string
	}{
		{2, "紧张和压力感"},
		{5, "事情发展不如预期"},
		{9, "失控感"},
		{15, "工作生活失衡"},
		{22, "情绪易怒"},
	}

	for _, item := range highScoreItemDescs {
		if item.index < len(userAnswers) && userAnswers[item.index] >= 4 {
			highScoreItems = append(highScoreItems, item.desc)
		}
	}

	// 改善建议
	suggestions = []string{
		"建立规律作息，保证充足睡眠",
		"每天进行30分钟有氧运动",
		"学习深呼吸等放松技巧",
		"与亲友分享感受，寻求支持",
		"合理安排时间，设置优先级",
	}

	return map[string]interface{}{
		"test_name":        "感知压力量表",
		"total_score":      totalScore,
		"max_score":        100,
		"level":            level,
		"title":            title,
		"analysis":         analysis,
		"high_score_items": highScoreItems,
		"suggestions":      suggestions,
	}
}
