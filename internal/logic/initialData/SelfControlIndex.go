package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// SelfControlIndex 自我控力指数评估初始化
func SelfControlIndex() {
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
		{Text: "我能够坚持完成即使是不喜欢的任务"},
		{Text: "我会为了长期目标而牺牲即时满足"},
		{Text: "我能控制自己的情绪，即使在压力下"},
		{Text: "我会制定详细的计划并按照计划执行"},
		{Text: "我能够抵制诱惑，专注于重要的事情"},
		{Text: "我会定期反思自己的行为并做出改进"},
		{Text: "我能在干扰环境中保持专注"},
		{Text: "我会为自己设定明确的目标"},
		{Text: "我能够控制自己的冲动行为"},
		{Text: "我会优先处理重要但不紧急的事情"},
		{Text: "我能够坚持健康的生活习惯"},
		{Text: "我会控制自己的消费行为，避免冲动购物"},
		{Text: "我能在面对挫折时保持积极态度"},
		{Text: "我会合理安排时间，避免拖延"},
		{Text: "我能够控制自己的饮食，保持健康体重"},
		{Text: "我会为未来的风险做准备"},
		{Text: "我能在生气时控制自己的言行"},
		{Text: "我会定期学习新技能，提升自己"},
		{Text: "我能够抵制社交媒体的诱惑，专注于工作"},
		{Text: "我会遵守自己设定的规则和承诺"},
		{Text: "我能够控制自己的娱乐时间，避免过度"},
		{Text: "我会为自己的行为负责"},
		{Text: "我能在面对困难时坚持不懈"},
		{Text: "我会优先考虑长期利益而非短期快乐"},
		{Text: "我能够控制自己的语言，避免说伤害他人的话"},
		{Text: "我会定期锻炼，保持身体健康"},
		{Text: "我能在面对批评时保持冷静并从中学习"},
		{Text: "我会制定预算并严格遵守"},
		{Text: "我能够控制自己的注意力，避免分心"},
		{Text: "我会为自己的错误道歉并努力改正"},
		{Text: "我能够抵制垃圾食品的诱惑"},
		{Text: "我会设定合理的边界，避免过度承诺"},
		{Text: "我能在面对压力时保持冷静"},
		{Text: "我会定期整理自己的环境，保持整洁"},
		{Text: "我能够控制自己的睡眠时间，保持规律"},
		{Text: "我会优先处理紧急且重要的事情"},
		{Text: "我能在面对诱惑时说'不'"},
		{Text: "我会为自己的成功设定更高的标准"},
		{Text: "我能够控制自己的嫉妒情绪"},
		{Text: "我会定期回顾自己的目标并调整计划"},
		{Text: "我能够抵制酒精或其他有害物质的诱惑"},
		{Text: "我会尊重他人的意见，即使与自己不同"},
		{Text: "我能在面对不确定性时保持冷静"},
		{Text: "我会定期储蓄，为未来做准备"},
		{Text: "我能够控制自己的竞争欲望，保持健康心态"},
		{Text: "我会优先考虑他人的需求，而非只关注自己"},
		{Text: "我能在面对失败时保持韧性"},
		{Text: "我会定期清理自己的数字空间，保持有序"},
		{Text: "我能够控制自己的完美主义倾向，避免过度追求"},
		{Text: "我会为自己的成就感到自豪，但不过度自满"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.SelfControlIndex, "自我控力指数评估", "自我控制量表(SCS)由Tangney、Baumeister和Boone于2004年编制，通过36个条目评估个体在日常行为规划、延迟满足、情绪调节、冲动抑制及习惯养成等方面的自我约束能力，研究表明高自控力与学业成就、心理健康及人际关系质量显著正相关，是自我控制研究的权威测量工具。", []ReferenceItem{{Authors: []string{"Tangney, J.P.", "Baumeister, R.F.", "Boone, A.L."}, Title: "High self-control predicts good adjustment, less pathology, better grades, and interpersonal success", Journal: "Journal of Personality", Year: 2004, DOI: "10.1111/j.0022-3506.2004.00263.x"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.SelfControlIndex,
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