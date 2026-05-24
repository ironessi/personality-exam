package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// MBTI 16型人格测试初始化
func MBTI() {
	ctx := context.Background()

	// 定义通用选项（7点计分，Tag A-G）
	options := []Option{
		{Tag: "A", Content: "非常不符合", Score: 1},
		{Tag: "B", Content: "不符合", Score: 2},
		{Tag: "C", Content: "有点不符合", Score: 3},
		{Tag: "D", Content: "中立", Score: 4},
		{Tag: "E", Content: "有点符合", Score: 5},
		{Tag: "F", Content: "符合", Score: 6},
		{Tag: "G", Content: "非常符合", Score: 7},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我更倾向于通过社交活动获取能量"},
		{Text: "我更关注事实和细节，而非抽象概念"},
		{Text: "我做决策时更依赖逻辑分析，而非情感因素"},
		{Text: "我喜欢计划和组织，倾向于有秩序的生活"},
		{Text: "我喜欢与他人合作，而非独立工作"},
		{Text: "我更相信经验，而非直觉"},
		{Text: "我更注重公平，而非和谐"},
		{Text: "我喜欢提前安排好事情，而非即兴发挥"},
		{Text: "我喜欢在大群体中社交"},
		{Text: "我更关注当前，而非未来"},
	}

	// 保存测试类型信息

	SaveTestType(ctx, int(consts.MBTI), "16型人格测试（MBTI）", "MBTI人格测试基于荣格心理类型理论，由Myers和Briggs于20世纪40年代开发，通过外倾/内倾、感觉/直觉、思考/情感、判断/感知四个维度的组合，将人格划分为16种类型，广泛应用于职业发展、团队协作、人际沟通与自我认知领域，是全球使用最广泛的人格评估工具之一。", []ReferenceItem{{Authors: []string{"Myers, I.B.", "McCaulley, M.H."}, Title: "Manual: A guide to the development and use of the Myers-Briggs Type Indicator", Journal: "Consulting Psychologists Press", Year: 1985}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.MBTI).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.MBTI,
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

		// 插入选项
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
