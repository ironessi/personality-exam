package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func SelfEfficacy() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "如果我努力去做的话，我总是能够解决问题的。"},
		{Text: "即使别人反对我，我仍有办法取得我所要的。"},
		{Text: "对我来说，坚持理想和达成目标是轻而易举的。"},
		{Text: "我自信能有效地应付任何突如其来的事情。"},
		{Text: "以我的才智，我定能应付意料之外的情况。"},
		{Text: "如果我付出必要的努力，我一定能解决大多数的难题。"},
		{Text: "我能冷静地面对困难，因为我信赖自己处理问题的能力。"},
		{Text: "面对一个难题时，我通常能找到几个解决方法。"},
		{Text: "有麻烦的时候，我通常能想到一些应付的方法。"},
		{Text: "无论什么事在我身上发生，我都能应付自如。"},
	}

	// 定义选项
	options := []Option{
		{Tag: "A", Content: "从不", Score: 1},
		{Tag: "B", Content: "很少", Score: 2},
		{Tag: "C", Content: "有时", Score: 3},
		{Tag: "D", Content: "经常", Score: 4},
		{Tag: "E", Content: "总是", Score: 5},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.SelfEfficacy), "一般自我效能感测试", "一般自我效能感量表(GSES)由Schwarzer和Jerusalem于1981年编制，共10个条目，评估个体面对各种困难情境时的总体自信程度与应对能力感，基于Bandura的自我效能理论，已被翻译成35种以上语言，在全球超过100个国家广泛应用于临床、教育与健康心理学研究。", []ReferenceItem{{Authors: []string{"Schwarzer, R.", "Jerusalem, M."}, Title: "Generalized Self-Efficacy scale", Journal: "Measures in health psychology: A user's portfolio", Year: 1995}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.SelfEfficacy).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.SelfEfficacy,
			Content:    question.Text,
			Sequence:   idx + 1,
		}).Insert()
		if err != nil {
			log.Printf("插入数据失败: %v\n", err)
			return
		}

		questionId, err := result.LastInsertId()
		if err != nil {
			log.Printf("获取插入ID失败: %v\n", err)
			return
		}

		// 插入选项 - 使用正确的JSON格式
		_, err = dao.Options.Ctx(ctx).Data(&do.Options{
			QuestionId: questionId,
			Content:    optionStr,
		}).Insert()
		if err != nil {
			log.Printf("插入选项失败: %v\n", err)
			return
		}
	}
}
