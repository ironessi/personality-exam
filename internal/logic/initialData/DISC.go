package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// DISC DISC性格测试初始化
func DISC() {
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
		{Text: "我喜欢掌控局面"},
		{Text: "我是一个行动果断的人"},
		{Text: "我喜欢设定具有挑战性的目标"},
		{Text: "我善于快速做出决策"},
		{Text: "我喜欢竞争"},
		{Text: "我是一个外向开朗的人"},
		{Text: "我喜欢与他人交往"},
		{Text: "我善于说服他人"},
		{Text: "我喜欢成为众人关注的焦点"},
		{Text: "我善于表达自己的想法"},
		{Text: "我是一个耐心细致的人"},
		{Text: "我喜欢稳定的生活和工作环境"},
		{Text: "我善于倾听他人的意见"},
		{Text: "我喜欢与他人和谐相处"},
		{Text: "我不喜欢冲突"},
		{Text: "我是一个注重细节的人"},
		{Text: "我喜欢按计划行事"},
		{Text: "我善于分析问题"},
		{Text: "我追求精确和准确性"},
		{Text: "我喜欢遵循规则和流程"},
		{Text: "我喜欢领导他人"},
		{Text: "我是一个结果导向的人"},
		{Text: "我喜欢冒险"},
		{Text: "我善于激励他人"},
		{Text: "我是一个充满活力的人"},
		{Text: "我喜欢社交活动"},
		{Text: "我善于建立人际关系"},
		{Text: "我是一个乐观向上的人"},
		{Text: "我喜欢帮助他人"},
		{Text: "我是一个可靠的人"},
		{Text: "我喜欢团队合作"},
		{Text: "我善于适应变化"},
		{Text: "我是一个谨慎的人"},
		{Text: "我喜欢解决复杂问题"},
		{Text: "我注重质量"},
		{Text: "我善于组织和规划"},
		{Text: "我是一个有条理的人"},
		{Text: "我喜欢独立工作"},
		{Text: "我善于时间管理"},
		{Text: "我追求效率"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.DISC), "DISC性格测试", "DISC性格测试基于Marston于1928年提出的情绪行为理论，将人格行为风格分为支配型(D)、影响型(I)、稳健型(S)、谨慎型(C)四种类型，通过识别个体的主导行为风格，帮助提升职场沟通效率、优化团队协作模式，是企业培训与领导力发展领域最广泛使用的行为评估工具之一。", []ReferenceItem{{Authors: []string{"Marston, W.M."}, Title: "Emotions of Normal People", Journal: "Kegan Paul Trench Trubner And Company", Year: 1928}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.DISC,
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