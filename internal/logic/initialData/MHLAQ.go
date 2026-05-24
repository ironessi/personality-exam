package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// MHLAQ 心理健康素养量表初始化
func MHLAQ() {
	ctx := context.Background()

	// 定义选项（5点计分）
	options := []Option{
		{Tag: "A", Content: "非常不同意", Score: 1},
		{Tag: "B", Content: "不同意", Score: 2},
		{Tag: "C", Content: "中立", Score: 3},
		{Tag: "D", Content: "同意", Score: 4},
		{Tag: "E", Content: "非常同意", Score: 5},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我知道如何识别常见的心理问题症状"},
		{Text: "我了解不同类型的心理治疗方法"},
		{Text: "我知道如何寻求专业的心理帮助"},
		{Text: "我认为心理健康和身体健康同样重要"},
		{Text: "我知道如何管理自己的情绪"},
		{Text: "我了解压力对心理健康的影响"},
		{Text: "我知道如何应对挫折和困难"},
		{Text: "我认为寻求心理帮助是一种勇敢的行为"},
		{Text: "我了解不同的心理健康服务资源"},
		{Text: "我知道如何帮助有心理问题的朋友或家人"},
		{Text: "我认为心理健康问题是可以预防的"},
		{Text: "我了解心理健康问题的常见诱因"},
		{Text: "我知道如何维持良好的心理健康"},
		{Text: "我认为心理健康问题是可以治疗的"},
		{Text: "我了解心理健康问题的去污名化"},
		{Text: "我知道如何在生活中保持平衡"},
		{Text: "我了解心理健康教育的重要性"},
		{Text: "我认为每个人都应该关注自己的心理健康"},
		{Text: "我知道如何识别心理健康危机"},
		{Text: "我了解心理健康问题对生活质量的影响"},
		{Text: "我知道如何在压力下保持冷静"},
		{Text: "我认为心理健康问题并不可耻"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.MHLAQ), "心理健康素养量表(MHLAQ)", "心理健康素养量表(MHLAQ)由Jorm等于1997年提出，评估个体识别心理疾病的能力、对心理健康知识的了解程度及寻求专业帮助的态度与意愿，有助于减少心理健康污名化、促进早期识别与干预，是公众心理健康教育效果评估的重要工具。", []ReferenceItem{{Authors: []string{"Jorm, A.F.", "Korten, A.E.", "Jacomb, P.A."}, Title: "Mental health literacy: A survey of the public's ability to recognise mental disorders", Journal: "Medical Journal of Australia", Year: 1997, DOI: "10.5694/j.1326-5377.1997.tb140923.x"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.MHLAQ,
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