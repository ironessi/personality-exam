package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// SDSAnxiety 焦虑症自测量表初始化
func SDSAnxiety() {
	ctx := context.Background()

	// 定义选项（SDS使用1-4分计分）
	options := []Option{
		{Tag: "A", Content: "很少", Score: 1},
		{Tag: "B", Content: "有时", Score: 2},
		{Tag: "C", Content: "经常", Score: 3},
		{Tag: "D", Content: "总是", Score: 4},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我觉得比平时容易紧张和着急"},
		{Text: "我无缘无故地感到害怕"},
		{Text: "我容易心里烦乱或觉得惊恐"},
		{Text: "我觉得我可能将要发疯"},
		{Text: "我觉得一切都很好，也不会发生什么不幸"},
		{Text: "我手脚发抖打颤"},
		{Text: "我因为头痛、头颈痛和背痛而苦恼"},
		{Text: "我感觉容易衰弱和疲乏"},
		{Text: "我觉得心平气和，并且容易安静坐着"},
		{Text: "我觉得心跳得很快"},
		{Text: "我因为一阵阵头晕而苦恼"},
		{Text: "我有晕倒发作，或觉得要晕倒似的"},
		{Text: "我呼气吸气都感到很容易"},
		{Text: "我手脚麻木和刺痛"},
		{Text: "我因为胃痛和消化不良而苦恼"},
		{Text: "我常常要小便"},
		{Text: "我的手脚常常是干燥温暖的"},
		{Text: "我脸红发热"},
		{Text: "我容易入睡并且一夜睡得很好"},
		{Text: "我做恶梦"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.SDSAnxiety, "焦虑症自测量表(SDS)", "焦虑自评量表(SAS)由Zung于1971年编制，包含20个条目，从主观感受、躯体症状和行为表现三个层面评估个体的焦虑症状及其严重程度，操作简便、耗时短，广泛应用于临床焦虑障碍筛查、心理健康流行病学调查及焦虑干预效果的追踪评估。", []ReferenceItem{{Authors: []string{"Zung, W.W.K."}, Title: "A rating instrument for anxiety disorders", Journal: "Psychosomatics", Year: 1971, DOI: "10.1016/S0033-3182(71)71479-0"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.SDSAnxiety,
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