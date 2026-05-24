package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// EMP 同理心通情量表初始化
func EMP() {
	ctx := context.Background()

	// 定义选项（5点计分）
	options := []Option{
		{Tag: "A", Content: "非常不符合", Score: 1},
		{Tag: "B", Content: "不符合", Score: 2},
		{Tag: "C", Content: "中立", Score: 3},
		{Tag: "D", Content: "符合", Score: 4},
		{Tag: "E", Content: "非常符合", Score: 5},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我能理解他人的感受"},
		{Text: "我善于察觉他人的情绪变化"},
		{Text: "我能够站在他人的角度思考问题"},
		{Text: "我对他人的痛苦感同身受"},
		{Text: "我善于倾听他人的诉说"},
		{Text: "我能够体会他人的处境"},
		{Text: "我对他人的需求很敏感"},
		{Text: "我能够理解他人的观点"},
		{Text: "我善于安抚情绪低落的人"},
		{Text: "我能够感受到他人的快乐"},
		{Text: "我对他人的困境感到关心"},
		{Text: "我善于理解他人的行为动机"},
		{Text: "我能够与他人建立情感连接"},
		{Text: "我对他人的情绪反应很敏感"},
		{Text: "我能够同情他人的遭遇"},
		{Text: "我善于与他人共情"},
		{Text: "我能够理解他人的压力"},
		{Text: "我对他人的感受很在意"},
		{Text: "我能够体会他人的孤独"},
		{Text: "我善于理解他人的恐惧"},
		{Text: "我对他人的喜悦感同身受"},
		{Text: "我能够理解他人的愤怒"},
		{Text: "我善于察觉他人的不安"},
		{Text: "我能够体会他人的失望"},
		{Text: "我对他人的痛苦感到难过"},
		{Text: "我能够理解他人的焦虑"},
		{Text: "我善于理解他人的期望"},
		{Text: "我能够与他人的情绪产生共鸣"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.EMP), "同理心通情量表（EMP）", "同理心通情量表(EMP)基于Davis于1980年提出的多维同理心理论，系统评估个体的情感共情（感受他人情绪）、认知共情（理解他人视角）和社会共情（理解社会情境）三个核心成分，广泛应用于临床心理学、社会工作、教育及人际关系研究领域。", []ReferenceItem{{Authors: []string{"Davis, M.H."}, Title: "Measuring individual differences in empathy: Evidence for a multidimensional approach", Journal: "Journal of Personality and Social Psychology", Year: 1983, DOI: "10.1037/0022-3514.44.1.113"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.EMP,
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