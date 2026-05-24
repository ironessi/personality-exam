package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func GeneralLifeSatisfaction() {
	ctx := context.Background()

	var questions = []Question{
		{
			Text: "你这个年纪最大的好处是什么？",
			Options: []Option{
				{Tag: "A", Content: "积极的答案", Score: 2},
				{Tag: "B", Content: "没有任何好处", Score: 0},
			},
		},
		{
			Text: "今后五年你打算做什么？你估计今后的生活会有什么变化？",
			Options: []Option{
				{Tag: "A", Content: "变好，或无变化", Score: 2},
				{Tag: "B", Content: "无法预料，'各种可能性都有'", Score: 1},
				{Tag: "C", Content: "变坏", Score: 0},
			},
		},
		{
			Text: "你现在生活中最重要的事情是什么？",
			Options: []Option{
				{Tag: "A", Content: "任何自身之外的事情，或令人愉快的对未来的解释", Score: 2},
				{Tag: "B", Content: "'维持现状'、保持健康或工作", Score: 1},
				{Tag: "C", Content: "摆脱现在的困境、或'目前什么重要的事情也没有'、或提起以往的经历", Score: 0},
			},
		},
		{
			Text: "与早期的生活相比，你现在是否幸福？",
			Options: []Option{
				{Tag: "A", Content: "现在是最幸福的时期，过去和现在同样幸福；或无法比较出何时更幸福", Score: 2},
				{Tag: "B", Content: "最近几年有些不如以前了", Score: 1},
				{Tag: "C", Content: "以前比现在好，目前是最糟糕的时期", Score: 0},
			},
		},
		{
			Text: "你是否曾担心人们期望你做的事你却不能胜任---你无法满足人们对你的要求？",
			Options: []Option{
				{Tag: "A", Content: "不曾担心", Score: 2},
				{Tag: "B", Content: "略有些担心", Score: 1},
				{Tag: "C", Content: "担心", Score: 0},
			},
		},
		{
			Text: "如果你想怎样就能怎样，那么你最喜欢生活在哪里（国家名）？",
			Options: []Option{
				{Tag: "A", Content: "目前所在地", Score: 2},
				{Tag: "B", Content: "任何其它地方", Score: 0},
			},
		},
		{
			Text: "你感到孤独的时间有多少？",
			Options: []Option{
				{Tag: "A", Content: "从未有过", Score: 2},
				{Tag: "B", Content: "有时", Score: 1},
				{Tag: "C", Content: "经常，十分频繁", Score: 0},
			},
		},
		{
			Text: "你感到生活无目的的时间有多少？",
			Options: []Option{
				{Tag: "A", Content: "从未有过", Score: 2},
				{Tag: "B", Content: "有时", Score: 1},
				{Tag: "C", Content: "经常，十分频繁", Score: 0},
			},
		},
		{
			Text: "你希望将来与好朋友在一起的时间更多一些还是自己独处的时间更多一些？",
			Options: []Option{
				{Tag: "A", Content: "现在这样很好", Score: 2},
				{Tag: "B", Content: "与好朋友在一起的时间更多一些", Score: 1},
				{Tag: "C", Content: "自己独处的时间更多一个", Score: 0},
			},
		},
		{
			Text: "你在目前的生活中发现多少不幸的事情？",
			Options: []Option{
				{Tag: "A", Content: "几乎没有", Score: 2},
				{Tag: "B", Content: "有一些", Score: 1},
				{Tag: "C", Content: "许多", Score: 0},
			},
		},
		{
			Text: "当你年迈之后，事情比原先想象得好还是不好？",
			Options: []Option{
				{Tag: "A", Content: "好", Score: 2},
				{Tag: "B", Content: "和预期的差不多", Score: 1},
				{Tag: "C", Content: "不好", Score: 0},
			},
		},
		{
			Text: "你对自己生活的满意程度如何？",
			Options: []Option{
				{Tag: "A", Content: "非常满意", Score: 2},
				{Tag: "B", Content: "相当满意", Score: 1},
				{Tag: "C", Content: "不太满意", Score: 0},
			},
		},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.GeneralLifeSatisfaction), "生活满意度", "生活满意度量表(LSIA)由Neugarten等于1961年编制，通过热情与活力、决心与毅力、目标与意义、自我概念及情绪基调五个维度，综合评估个体对自身生活状态的总体主观满意程度，是老年心理学与积极心理学领域的经典测量工具，适用于成人及老年群体的幸福感研究。", []ReferenceItem{{Authors: []string{"Neugarten, B.L.", "Havighurst, R.J.", "Tobin, S.S."}, Title: "The measurement of life satisfaction", Journal: "Journal of Gerontology", Year: 1961}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.GeneralLifeSatisfaction).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		// 插入问题
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.GeneralLifeSatisfaction,
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

		// 将选项转换为JSON
		optionStr := optionsToJSON(question.Options)

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
