package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// ObedienceTest 服从性测试初始化
func ObedienceTest() {
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
		{Text: "我通常会遵守权威人士的指示，即使我不完全同意"},
		{Text: "我认为服从规则和法律是公民的责任"},
		{Text: "我会为了避免冲突而遵循他人的建议"},
		{Text: "我相信专业人士的判断通常比普通人更准确"},
		{Text: "我会优先考虑集体利益而非个人偏好"},
		{Text: "我认为尊重权威是一种美德"},
		{Text: "我会执行上级的命令，即使我认为这不是最佳选择"},
		{Text: "我重视传统价值观和社会规范"},
		{Text: "我会为了维持和谐关系而妥协自己的立场"},
		{Text: "我相信大多数情况下，服从是明智的选择"},
		{Text: "我会遵循既定的流程和程序"},
		{Text: "我认为质疑权威可能会导致混乱"},
		{Text: "我会按照指示完成任务，即使这需要额外的努力"},
		{Text: "我重视团队的一致性和凝聚力"},
		{Text: "我会避免挑战现状，除非有充分的理由"},
		{Text: "我认为服从是建立信任和尊重的基础"},
		{Text: "我会遵循长辈的建议，即使我有不同的想法"},
		{Text: "我相信规则的存在是为了保护每个人的利益"},
		{Text: "我会为了维护秩序而做出个人牺牲"},
		{Text: "我认为服从并不意味着缺乏独立思考能力"},
		{Text: "我会执行群体的决定，即使我个人持保留意见"},
		{Text: "我重视专业知识和经验"},
		{Text: "我会遵循行业标准和最佳实践"},
		{Text: "我认为服从是一种成熟的表现"},
		{Text: "我会为了避免麻烦而遵循他人的要求"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.ObedienceTest), "服从性测试", "服从性测试基于Milgram于1963年开展的经典服从权威实验理论框架，评估个体在权威压力与社会规范下的服从倾向与自主意志强度，探索人际关系中的权力动态、个体边界设定能力及道德判断与服从行为之间的心理冲突，适用于自我认知与人际关系模式的深度探索。", []ReferenceItem{{Authors: []string{"Milgram, S."}, Title: "Behavioral study of obedience", Journal: "Journal of Abnormal and Social Psychology", Year: 1963, DOI: "10.1037/h0040525"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.ObedienceTest,
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