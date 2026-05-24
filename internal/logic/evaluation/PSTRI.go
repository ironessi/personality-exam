package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

// PSTRIResult PSTRI压力知觉测试量表结果
type PSTRIResult struct {
	TotalScore  int      // 总分
	Level       string   // 等级
	Title       string   // 标题
	Analysis    string   // 分析
	Suggestions []string // 建议
}

// PSTRI PSTRI压力知觉测试量表评分分析
func (e *Evaluation) PSTRI(data []v1.AnswerRecord) (interface{}, error) {
	// 计算总分
	totalScore := 0

	// 正向计分题目：1, 2, 3, 6, 8, 9, 10, 12, 13
	// 反向计分题目：4, 5, 7, 11, 14, 15

	// 简单的计分逻辑
	for _, record := range data {
		score := int(record.Score)
		totalScore += score
	}

	// 确定等级和分析
	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	title = "PSTRI压力知觉测试量表结果"
	analysis = fmt.Sprintf("您的PSTRI测试得分为：%d分\n\n压力知觉水平：", totalScore)

	// 根据总分确定压力知觉水平
	if totalScore <= 10 {
		level = "低压力"
		analysis += "低压力水平\n\n您目前感受到的压力水平较低，能够较好地应对生活中的各种挑战。"
		suggestions = []string{
			"继续保持良好的生活习惯和应对方式",
			"定期关注自己的情绪变化",
			"保持积极的心态，享受生活",
			"适当参加一些有意义的活动，丰富自己的生活",
			"与家人和朋友保持良好的沟通",
		}
	} else if totalScore <= 20 {
		level = "中等压力"
		analysis += "中等压力水平\n\n您目前感受到的压力水平适中，需要适当调整自己的生活方式和应对策略。"
		suggestions = []string{
			"学习一些压力管理技巧，如深呼吸、冥想等",
			"合理安排工作和休息时间，避免过度劳累",
			"保持良好的睡眠质量，每天保证充足的睡眠时间",
			"适当进行体育锻炼，释放压力",
			"与家人和朋友分享自己的感受，寻求支持",
		}
	} else if totalScore <= 30 {
		level = "高压力"
		analysis += "高压力水平\n\n您目前感受到的压力水平较高，需要积极采取措施来缓解压力。"
		suggestions = []string{
			"学习有效的压力管理方法，如时间管理、情绪调节等",
			"适当降低自己的期望和要求，避免过度追求完美",
			"增加休息时间，保证充足的睡眠",
			"寻求专业心理咨询师的帮助，学习应对压力的技巧",
			"培养一些兴趣爱好，转移注意力",
		}
	} else {
		level = "极高压力"
		analysis += "极高压力水平\n\n您目前感受到的压力水平极高，需要立即采取措施来缓解压力，避免对身心健康造成严重影响。"
		suggestions = []string{
			"立即寻求专业心理咨询师或医生的帮助",
			"暂时减少工作和生活压力，给自己充足的休息时间",
			"学习放松技巧，如渐进式肌肉放松、冥想等",
			"与家人和朋友保持密切联系，寻求他们的支持和帮助",
			"考虑调整自己的生活方式和工作环境，减少压力源",
		}
	}

	return map[string]interface{}{
		"test_name":   "PSTRI压力知觉测试量表",
		"total_score": totalScore,
		"level":       level,
		"title":       title,
		"analysis":    analysis,
		"suggestions": suggestions,
	}, nil
}