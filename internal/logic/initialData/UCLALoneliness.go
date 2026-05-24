package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func UCLALoneliness() {
	ctx := context.Background()

	// 正向计分选项
	optionTrue := []Option{
		{Tag: "A", Content: "从不", Score: 4},
		{Tag: "B", Content: "很少", Score: 3},
		{Tag: "C", Content: "有时", Score: 2},
		{Tag: "D", Content: "一直", Score: 1},
	}

	// 反向计分选项
	optionFalse := []Option{
		{Tag: "A", Content: "从不", Score: 1},
		{Tag: "B", Content: "很少", Score: 2},
		{Tag: "C", Content: "有时", Score: 3},
		{Tag: "D", Content: "一直", Score: 4},
	}

	// 转换为JSON
	optionTrueStr := optionsToJSON(optionTrue)
	optionFalseStr := optionsToJSON(optionFalse)

	var questions = []Question{
		// 反向计分问题（积极的社会联系感受）
		{Text: "你常感到与周围人的关系和谐吗？", OptionText: optionTrueStr},
		{Text: "你常感到属于朋友们中的一员吗？", OptionText: optionTrueStr},
		{Text: "你常感到与周围的人有许多共同点吗？", OptionText: optionTrueStr},
		{Text: "你常感到想要与人来往、结交朋友吗？", OptionText: optionTrueStr},
		{Text: "你常感到与人亲近吗？", OptionText: optionTrueStr},
		{Text: "你常感到当你愿意时就能找到伙伴吗？", OptionText: optionTrueStr},
		{Text: "你常感到有人真正了解你吗？", OptionText: optionTrueStr},
		{Text: "你常感到有人愿意与你交谈吗？", OptionText: optionTrueStr},
		{Text: "你常感到有人值得你信赖吗？", OptionText: optionTrueStr},

		// 正向计分问题（消极的孤独感受）
		{Text: "你常感到缺少伙伴吗？", OptionText: optionFalseStr},
		{Text: "你常感到没人可以信赖吗？", OptionText: optionFalseStr},
		{Text: "你常感到寂寞吗？", OptionText: optionFalseStr},
		{Text: "你常感到与任何人都不亲密吗？", OptionText: optionFalseStr},
		{Text: "你常感到你的兴趣与想法与周围的人不一样吗？", OptionText: optionFalseStr},
		{Text: "你常感到被人冷落吗？", OptionText: optionFalseStr},
		{Text: "你常感到你与别人来往毫无意义吗？", OptionText: optionFalseStr},
		{Text: "你常感到没有人很了解你吗？", OptionText: optionFalseStr},
		{Text: "你常感到与别人隔开了吗？", OptionText: optionFalseStr},
		{Text: "你常感到羞怯吗？", OptionText: optionFalseStr},
		{Text: "你常感到人们围着你但并不关心你吗？", OptionText: optionFalseStr},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.UCLALoneliness), "孤独感", "UCLA孤独感量表由Russell、Peplau和Ferguson于1978年编制，第三版包含20个条目，通过评估个体对社交互动质量的满意程度、归属感的缺失以及与他人的心理距离，全面测量主观孤独感体验，是孤独研究领域引用最广泛的标准化工具，适用于青少年至老年各年龄段。", []ReferenceItem{{Authors: []string{"Russell, D.", "Peplau, L.A.", "Ferguson, M.L."}, Title: "Developing a measure of loneliness", Journal: "Journal of Personality Assessment", Year: 1978, DOI: "10.1207/s15327752jpa4203_1"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.UCLALoneliness).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.UCLALoneliness,
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

		// 使用正确的选项文本
		_, err = dao.Options.Ctx(ctx).Data(&do.Options{
			QuestionId: questionId,
			Content:    question.OptionText,
		}).Insert()
		if err != nil {
			log.Printf("插入选项失败 [%d]: %v\n", idx+1, err)
			break
		}
	}
}
