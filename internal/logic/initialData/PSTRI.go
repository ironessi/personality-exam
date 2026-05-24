package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// PSTRI PSTRI压力知觉测试量表初始化
func PSTRI() {
	ctx := context.Background()

	// 正向计分选项 (0-4)
	optionsForward := []Option{
		{Tag: "A", Content: "从不", Score: 0},
		{Tag: "B", Content: "很少", Score: 1},
		{Tag: "C", Content: "有时", Score: 2},
		{Tag: "D", Content: "经常", Score: 3},
		{Tag: "E", Content: "总是", Score: 4},
	}
	optionForwardStr := optionsToJSON(optionsForward)

	// 反向计分选项 (4-0)
	optionsReverse := []Option{
		{Tag: "A", Content: "从不", Score: 4},
		{Tag: "B", Content: "很少", Score: 3},
		{Tag: "C", Content: "有时", Score: 2},
		{Tag: "D", Content: "经常", Score: 1},
		{Tag: "E", Content: "总是", Score: 0},
	}
	optionReverseStr := optionsToJSON(optionsReverse)

	var questions = []Question{
		{Text: "我因为意外的事情而感到不安", OptionText: optionForwardStr},
		{Text: "我感到无法控制生活中的重要事情", OptionText: optionForwardStr},
		{Text: "我感到紧张和压力", OptionText: optionForwardStr},
		{Text: "我对处理个人问题的能力有信心", OptionText: optionReverseStr},
		{Text: "我觉得事情按照我的意愿发展", OptionText: optionReverseStr},
		{Text: "我发现自己在处理太多的事情", OptionText: optionForwardStr},
		{Text: "我能应对生活中的所有挑战", OptionText: optionReverseStr},
		{Text: "我觉得事情完全失控", OptionText: optionForwardStr},
		{Text: "我感到愤怒，因为太多事情超出我的控制", OptionText: optionForwardStr},
		{Text: "我感到困难堆积如山，无法克服", OptionText: optionForwardStr},
		{Text: "我觉得自己能有效地处理生活中的问题", OptionText: optionReverseStr},
		{Text: "我因为事情超出我的控制而感到沮丧", OptionText: optionForwardStr},
		{Text: "我对生活中的压力感到不堪重负", OptionText: optionForwardStr},
		{Text: "我感到有能力处理自己的情绪", OptionText: optionReverseStr},
		{Text: "我觉得事情进展顺利", OptionText: optionReverseStr},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.PSTRI, "PSTRI压力知觉测试量表", "压力知觉量表(PSS)由Cohen等于1983年编制，通过14个条目评估个体对生活情境中不可预测性、不可控制性及超负荷感的主观认知评价，不依赖特定压力事件，而是测量个体对整体生活压力的感知强度，是心理压力研究领域最广泛使用的标准化自评工具。", []ReferenceItem{{Authors: []string{"Cohen, S.", "Kamarck, T.", "Mermelstein, R."}, Title: "A global measure of perceived stress", Journal: "Journal of Health and Social Behavior", Year: 1983, DOI: "10.2307/2136404"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.PSTRI).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.PSTRI,
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

		// 插入选项 (使用每个问题指定的 OptionText)
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
