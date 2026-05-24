package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func Panas() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "害怕的"},
		{Text: "内疚的"},
		{Text: "紧张的"},
		{Text: "恼怒的"},
		{Text: "难过的"},
		{Text: "感激的"},
		{Text: "活跃的"},
		{Text: "兴奋的"},
		{Text: "自豪的"},
		{Text: "欣喜的"},
		{Text: "惊恐的"},
		{Text: "羞愧的"},
		{Text: "易怒的"},
		{Text: "快乐的"},
		{Text: "精神充沛的"},
		{Text: "兴高采烈的"},
		{Text: "充满热情的"},
		{Text: "战战兢兢的"},
	}

	// 定义选项
	options := []Option{
		{Tag: "A", Content: "完全没有", Score: 1},
		{Tag: "B", Content: "有一点", Score: 2},
		{Tag: "C", Content: "中等", Score: 3},
		{Tag: "D", Content: "相当多", Score: 4},
		{Tag: "E", Content: "非常多", Score: 5},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.Panas), "积极消极情绪量表", "积极消极情绪量表(PANAS)由Watson、Clark和Tellegen于1988年编制，包含20个情绪形容词，分别测量积极情感(PA)与消极情感(NA)两个相对独立的维度，通过评估特定时间段内各情绪词的体验频率，反映个体的主观情感平衡状态，是情绪与心理健康研究领域最广泛使用的标准化工具之一。", []ReferenceItem{{Authors: []string{"Watson, D.", "Clark, L.A.", "Tellegen, A."}, Title: "Development and validation of brief measures of positive and negative affect", Journal: "Journal of Personality and Social Psychology", Year: 1988, DOI: "10.1037/0022-3514.54.6.1063"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.Panas).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.Panas,
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
