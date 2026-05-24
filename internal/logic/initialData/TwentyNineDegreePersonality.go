package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// TwentyNineDegreePersonality 29度人格测试初始化
func TwentyNineDegreePersonality() {
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
		{Text: "面对阶段收尾时，我更感到焦虑而非从容"},
		{Text: "掌握技能后，我渴望突破新领域"},
		{Text: "我常在'坚守固有模式'与'尝试改变'间摇摆"},
		{Text: "我对自己的成长有清晰的阶段性规划"},
		{Text: "我善于在结束中寻找新的开始"},
		{Text: "我对未知的变化感到兴奋而非恐惧"},
		{Text: "我能很好地平衡过去经验与未来可能性"},
		{Text: "我认为成长是一个持续的过程，没有终点"},
		{Text: "我善于在关键时刻做出决策"},
		{Text: "我对自己的能力有清晰的认知"},
		{Text: "我愿意接受挑战，即使可能失败"},
		{Text: "我善于从经验中学习并应用到新情境"},
		{Text: "我对生活中的转折点有敏锐的感知"},
		{Text: "我能够很好地适应变化"},
		{Text: "我对自己的人生方向有明确的把握"},
		{Text: "我善于在复杂情境中保持清晰的思维"},
		{Text: "我愿意尝试新的方法和途径"},
		{Text: "我对自己的成长进度感到满意"},
		{Text: "我善于在压力下保持冷静"},
		{Text: "我对未来充满信心"},
		{Text: "我善于在不同角色间转换"},
		{Text: "我能够很好地处理生活中的矛盾"},
		{Text: "我对自己的优缺点有客观的认识"},
		{Text: "我愿意接受他人的反馈和建议"},
		{Text: "我善于在团队中发挥自己的作用"},
		{Text: "我对自己的人生有掌控感"},
		{Text: "我善于在变化中找到机会"},
		{Text: "我对自己的未来有明确的目标"},
		{Text: "我能够很好地平衡工作和生活"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.TwentyNineDegreePersonality, "29度人格测试", "29度人格测试根植于占星学的临界度数理论，聚焦黄道十二宫每宫最后一度的特殊能量场，以特质成熟度为核心评估维度，探索个体性格中兼具终结性与开创性的独特张力，帮助用户理解自身在人生关键转折点上的心理特征与潜在能量，适合对自我探索有深度兴趣的用户。", []ReferenceItem{{Authors: []string{"Jung, C.G."}, Title: "Psychological Types", Journal: "Princeton University Press", Year: 1971}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.TwentyNineDegreePersonality,
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