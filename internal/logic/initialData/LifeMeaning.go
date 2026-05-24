package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func LifeMeaning() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "我正在寻觅我人生的一个目的或使命"},
		{Text: "我的生活没有明确的目的"},
		{Text: "我正在寻找自己生活的意义"},
		{Text: "我明白自己生活的意义"},
		{Text: "我正在寻觅让我感觉自己生活饶有意义的东西"},
		{Text: "我总在尝试找寻自己生活的目的"},
		{Text: "我的生活有一个清晰的方向"},
		{Text: "我知道什么东西能使自己的生活有意义"},
		{Text: "我已经发现一个让自己满意的生活目的"},
	}

	// 正向计分选项
	options := []Option{
		{Tag: "A", Content: "完全不符合", Score: 1},
		{Tag: "B", Content: "比较不符合", Score: 2},
		{Tag: "C", Content: "有点不符合", Score: 3},
		{Tag: "D", Content: "不确定", Score: 4},
		{Tag: "E", Content: "有点符合", Score: 5},
		{Tag: "F", Content: "比较符合", Score: 6},
		{Tag: "G", Content: "完全符合", Score: 7},
	}

	// 反向计分选项（第2题使用）
	reverseOptions := []Option{
		{Tag: "A", Content: "完全不符合", Score: 7},
		{Tag: "B", Content: "比较不符合", Score: 6},
		{Tag: "C", Content: "有点不符合", Score: 5},
		{Tag: "D", Content: "不确定", Score: 4},
		{Tag: "E", Content: "有点符合", Score: 3},
		{Tag: "F", Content: "比较符合", Score: 2},
		{Tag: "G", Content: "完全符合", Score: 1},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)
	reverseOptionStr := optionsToJSON(reverseOptions)

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.LifeMeaning), "生命意义感测试", "生命意义感量表(MLQ)由Steger等于2006年编制，包含意义存在与意义寻求两个独立维度，前者评估个体当前感受到的生活目的与意义感，后者评估主动探索生命意义的动机强度，广泛应用于积极心理学、存在主义心理治疗及生命教育领域的研究与实践。", []ReferenceItem{{Authors: []string{"Steger, M.F.", "Frazier, P.", "Oishi, S.", "Kaler, M."}, Title: "The meaning in life questionnaire", Journal: "Journal of Counseling Psychology", Year: 2006, DOI: "10.1037/0022-0167.53.1.80"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.LifeMeaning).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.LifeMeaning,
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

		// 第2题使用反向计分选项，其他使用正向计分选项
		var currentOptionStr string
		if idx == 1 { // 第2题（索引为1）
			currentOptionStr = reverseOptionStr
		} else {
			currentOptionStr = optionStr
		}

		_, err = dao.Options.Ctx(ctx).Data(&do.Options{
			QuestionId: questionId,
			Content:    currentOptionStr,
		}).Insert()
		if err != nil {
			log.Printf("插入选项失败: %v\n", err)
			return
		}
	}
}
