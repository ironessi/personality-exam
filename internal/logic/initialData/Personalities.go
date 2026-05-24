package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// Personalities Personalities人格测试初始化
func Personalities() {
	ctx := context.Background()

	// 定义通用选项（添加Tags）
	options := []Option{
		{Tag: "A", Content: "非常不符合", Score: 1},
		{Tag: "B", Content: "不太符合", Score: 2},
		{Tag: "C", Content: "比较符合", Score: 3},
		{Tag: "D", Content: "非常符合", Score: 4},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我通常是一个外向的人，喜欢与人交往"},
		{Text: "我喜欢尝试新事物和新体验"},
		{Text: "我做事通常很有计划性和组织性"},
		{Text: "我对他人的感受很敏感"},
		{Text: "我喜欢独处，享受安静的时光"},
		{Text: "我对细节很敏感，注重准确性"},
		{Text: "我喜欢思考抽象的概念和理论"},
		{Text: "我做决策时更依赖逻辑分析"},
		{Text: "我喜欢社交活动，喜欢认识新朋友"},
		{Text: "我喜欢按照计划行事，不喜欢意外变化"},
		{Text: "我富有想象力和创造力"},
		{Text: "我做决策时更考虑他人的感受"},
		{Text: "我喜欢独立工作，不依赖他人"},
		{Text: "我注重实际，喜欢解决现实问题"},
		{Text: "我喜欢探索新的想法和可能性"},
		{Text: "我通常很客观，不会被情绪影响判断"},
		{Text: "我喜欢与他人合作，共同完成任务"},
		{Text: "我做事认真负责，总是按时完成任务"},
		{Text: "我对艺术和美学很感兴趣"},
		{Text: "我善于理解和同情他人"},
	}

	// 插入或更新测试类型
	SaveTestType(ctx, int(consts.Personalities), "Personalities人格测试", "大五人格测试(Big Five)基于Costa和McCrae于1992年完善的五因素模型，通过开放性、尽责性、外向性、宜人性、神经质五个维度全面描述人格结构，具有跨文化稳定性，是当代人格心理学研究的主流框架，广泛应用于临床评估、职业选拔与人际关系研究。", []ReferenceItem{{Authors: []string{"Costa, P.T.", "McCrae, R.R."}, Title: "Revised NEO Personality Inventory (NEO-PI-R) and NEO Five-Factor Inventory (NEO-FFI) professional manual", Journal: "Psychological Assessment Resources", Year: 1992}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.Personalities).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.Personalities,
			Content:    question.Text,
			Sequence:   idx + 1,
		}).Insert()
		if err != nil {
			log.Printf("插入问题失败 [%d]: %v\n", idx+1, err)
			break
		}

		questionId, err := result.LastInsertId()
		if err != nil {
			log.Printf("获取问题ID失败 [%d]: %v\n", idx+1, err)
			break
		}

		// 插入选项 (使用通用选项字符串)
		_, err = dao.Options.Ctx(ctx).Data(&do.Options{
			QuestionId: questionId,
			Content:    optionStr,
		}).Insert()
		if err != nil {
			log.Printf("插入选项失败 [%d]: %v\n", idx+1, err)
			break
		}
	}
}
