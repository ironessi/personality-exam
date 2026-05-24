package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
)

func (e *Evaluation) LifeMeaning(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	presenceScore := 0 // 存在意义维度分
	searchScore := 0   // 寻找意义维度分

	// 存在意义维度题目：第2,4,7,8,9题（索引1,3,6,7,8）
	presenceItems := []int{1, 3, 6, 7, 8}
	// 寻找意义维度题目：第1,3,5,6题（索引0,2,4,5）
	searchItems := []int{0, 2, 4, 5}

	for i, record := range data {
		score := int(record.Score)
		totalScore += score

		// 计算存在意义维度分
		for _, idx := range presenceItems {
			if i == idx {
				presenceScore += score
				break
			}
		}

		// 计算寻找意义维度分
		for _, idx := range searchItems {
			if i == idx {
				searchScore += score
				break
			}
		}
	}

	level := ""
	title := ""
	analysis := ""
	var suggestions []string

	// 生命意义感量表分级标准
	if totalScore >= 50 {
		level = "高意义感"
		title = "生命意义感强烈"
		analysis = fmt.Sprintf("您的生命意义感总分为%d分（满分63分）。其中存在意义维度%d分（满分35分），寻找意义维度%d分（满分28分）。您对生命有强烈的意义感和目的感，生活方向明确，能够从日常生活中找到价值和意义。",
			totalScore, presenceScore, searchScore)
		suggestions = []string{
			"继续保持积极的生活态度和目标导向",
			"尝试帮助他人找到生命意义和目的",
			"设定更具挑战性但可实现的人生目标",
			"记录和反思生命中有意义的时刻",
			"参与志愿服务，传播正能量",
		}
	} else if totalScore >= 36 {
		level = "中等意义感"
		title = "生命意义感适中"
		analysis = fmt.Sprintf("您的生命意义感总分为%d分（满分63分）。其中存在意义维度%d分（满分35分），寻找意义维度%d分（满分28分）。您对生命意义有一定感受，但可能在某些方面感到迷茫或不确定。",
			totalScore, presenceScore, searchScore)
		suggestions = []string{
			"多思考个人价值观和人生目标",
			"设定明确的短期和长期目标",
			"参与有意义的活动和社交团体",
			"练习感恩日记，记录积极体验",
			"寻求人生导师或参加成长工作坊",
		}
	} else {
		level = "低意义感"
		title = "生命意义感较弱"
		analysis = fmt.Sprintf("您的生命意义感总分为%d分（满分63分）。其中存在意义维度%d分（满分35分），寻找意义维度%d分（满分28分）。您可能正在经历生命意义的迷茫期，对生活方向和目的感到不确定。",
			totalScore, presenceScore, searchScore)
		suggestions = []string{
			"寻求心理咨询或生命教练的帮助",
			"探索不同的兴趣爱好和潜能",
			"练习正念冥想，关注当下体验",
			"阅读关于生命意义和人生哲学的书籍",
			"参与社区活动，建立有意义的连接",
		}
	}

	return map[string]interface{}{
		"test_name":      "生命意义感测试",
		"total_score":    totalScore,
		"presence_score": presenceScore,
		"search_score":   searchScore,
		"max_score":      63,
		"level":          level,
		"title":          title,
		"analysis":       analysis,
		"suggestions":    suggestions,
	}, nil
}
