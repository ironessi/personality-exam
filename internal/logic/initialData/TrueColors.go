package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// TrueColors 色彩性格测试初始化
func TrueColors() {
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
		{Text: "我重视和谐的人际关系"},
		{Text: "我喜欢有明确的规则和结构"},
		{Text: "我善于分析问题和寻找逻辑解决方案"},
		{Text: "我喜欢尝试新的体验和冒险"},
		{Text: "我关心他人的感受和需求"},
		{Text: "我注重细节和准确性"},
		{Text: "我喜欢挑战传统思维"},
		{Text: "我富有创造力和想象力"},
		{Text: "我善于倾听和理解他人"},
		{Text: "我喜欢按计划行事"},
		{Text: "我注重事实和数据"},
		{Text: "我喜欢 spontaneity 和灵活性"},
		{Text: "我重视团队合作和集体利益"},
		{Text: "我喜欢整洁有序的环境"},
		{Text: "我善于解决复杂问题"},
		{Text: "我喜欢成为众人关注的焦点"},
		{Text: "我富有同情心和同理心"},
		{Text: "我遵守承诺和责任"},
		{Text: "我喜欢独立思考和工作"},
		{Text: "我喜欢社交活动和与人交往"},
		{Text: "我重视和谐的家庭关系"},
		{Text: "我喜欢完成任务并看到成果"},
		{Text: "我喜欢学习新的知识和技能"},
		{Text: "我喜欢充满活力和变化的生活"},
		{Text: "我善于调解冲突和解决争端"},
		{Text: "我注重效率和实用性"},
		{Text: "我喜欢质疑权威和传统"},
		{Text: "我喜欢表达自己的情感和想法"},
		{Text: "我重视公平和正义"},
		{Text: "我喜欢设定目标并努力实现"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.TrueColors, "色彩性格测试（True Colors）", "色彩性格测试(True Colors)由Don Lowry于1978年基于Keirsey气质理论开发，将人格特质映射为蓝色（共情与理想主义）、金色（秩序与责任感）、绿色（逻辑与创新）、橙色（活力与适应性）四种色彩维度，以直观易懂的方式帮助个体了解自身性格优势，广泛应用于教育、团队建设与职业发展领域。", []ReferenceItem{{Authors: []string{"Lowry, D."}, Title: "True Colors: A Self-Discovery Activity for Understanding Others", Journal: "True Colors International", Year: 1978}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.TrueColors,
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