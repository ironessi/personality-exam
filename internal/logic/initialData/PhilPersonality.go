package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// PhilPersonality 菲尔人格十项测试题初始化
func PhilPersonality() {
	ctx := context.Background()

	// 定义通用选项（添加Tags）
	options := []Option{
		{Tag: "A", Content: "非常不符合", Score: 1},
		{Tag: "B", Content: "不太符合", Score: 2},
		{Tag: "C", Content: "比较符合", Score: 3},
		{Tag: "D", Content: "非常符合", Score: 4},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我喜欢参加社交活动，认识新朋友"},
		{Text: "我倾向于信任他人，认为大多数人是善良的"},
		{Text: "我做事认真负责，总是尽力完成任务"},
		{Text: "我情绪稳定，很少大起大落"},
		{Text: "我喜欢尝试新事物，对未知充满好奇"},
		{Text: "我喜欢独立思考，不随波逐流"},
		{Text: "我对他人的感受很敏感，善于理解他人"},
		{Text: "我喜欢有计划地安排生活，不喜欢混乱"},
		{Text: "我面对压力时能保持冷静，有效应对"},
		{Text: "我对艺术、音乐和文学有浓厚兴趣"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.PhilPersonality, "菲尔人格十项测试题", "菲尔人格测试由美国著名心理学家Dr. Phil McGraw编制，通过评估个人在日常生活中的习惯选择与行为反应模式，揭示个体潜在的性格倾向与心理特征，以通俗易懂的方式帮助大众了解自身人格，广泛应用于心理健康科普教育与自我认知提升领域。", []ReferenceItem{{Authors: []string{"McGraw, P.C."}, Title: "Self Matters: Creating Your Life from the Inside Out", Journal: "Free Press", Year: 2001}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		// 插入问题
		result, err := dao.Questions.Ctx(ctx).Data(do.Questions{
			TestTypeId: consts.PhilPersonality,
			Content:    question.Text, // 使用 Content 字段映射到数据库的 content 列
			Sequence:   idx + 1,
		}).Insert()

		if err != nil {
			g.Log().Error(ctx, "插入菲尔人格问题失败:", err)
			break
		}

		questionId, err := result.LastInsertId()
		if err != nil {
			g.Log().Error(ctx, "获取菲尔人格问题ID失败:", err)
			break
		}

		// 插入选项 (使用通用选项字符串)
		_, err = dao.Options.Ctx(ctx).Data(do.Options{
			QuestionId: questionId,
			Content:    optionStr,
		}).Insert()
		if err != nil {
			g.Log().Error(ctx, "插入菲尔人格选项失败:", err)
			break
		}
	}
	g.Log().Info(ctx, "菲尔人格十项测试题初始化完成")
}
