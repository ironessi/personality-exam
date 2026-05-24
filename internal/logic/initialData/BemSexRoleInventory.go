package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// BemSexRoleInventory 贝姆性别角色测试量表初始化
func BemSexRoleInventory() {
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
		// 男性特质
		{Text: "我是一个独立的人"},
		{Text: "我喜欢竞争和挑战"},
		{Text: "我有领导能力"},
		{Text: "我是一个果断的人"},
		{Text: "我喜欢冒险"},
		{Text: "我是一个自信的人"},
		{Text: "我有抱负和野心"},
		{Text: "我是一个坚强的人"},
		{Text: "我喜欢运动和体力活动"},
		{Text: "我是一个理性和逻辑的思考者"},
		
		// 女性特质
		{Text: "我是一个温柔的人"},
		{Text: "我喜欢照顾他人"},
		{Text: "我有同理心和 compassion"},
		{Text: "我是一个敏感的人"},
		{Text: "我喜欢艺术和音乐"},
		{Text: "我是一个有教养的人"},
		{Text: "我喜欢烹饪和家务"},
		{Text: "我是一个感性的人"},
		{Text: "我喜欢与他人分享感受"},
		{Text: "我是一个有耐心的人"},
		
		// 男性特质（续）
		{Text: "我是一个勇敢的人"},
		{Text: "我喜欢解决问题"},
		{Text: "我是一个有责任感的人"},
		{Text: "我喜欢掌控局面"},
		{Text: "我是一个有主见的人"},
		{Text: "我喜欢技术和机械"},
		{Text: "我是一个自信的领导者"},
		{Text: "我喜欢制定计划和目标"},
		{Text: "我是一个有决断力的人"},
		{Text: "我喜欢独立工作"},
		
		// 女性特质（续）
		{Text: "我是一个善解人意的人"},
		{Text: "我喜欢帮助他人"},
		{Text: "我是一个有创造力的人"},
		{Text: "我喜欢装饰和美化环境"},
		{Text: "我是一个有直觉的人"},
		{Text: "我喜欢与他人建立亲密关系"},
		{Text: "我是一个有爱心的人"},
		{Text: "我喜欢表达自己的情感"},
		{Text: "我是一个有耐心的倾听者"},
		{Text: "我喜欢团队合作和协作"},
		
		// 中性特质
		{Text: "我是一个诚实的人"},
		{Text: "我喜欢学习新事物"},
		{Text: "我是一个有责任感的人"},
		{Text: "我喜欢旅行和探索"},
		{Text: "我是一个有幽默感的人"},
		{Text: "我喜欢阅读和学习"},
		{Text: "我是一个有条理的人"},
		{Text: "我喜欢社交活动"},
		{Text: "我是一个有正义感的人"},
		{Text: "我喜欢享受生活"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.BemSexRoleInventory), "贝姆性别角色测试量表", "贝姆性别角色量表(BSRI)由Sandra Bem于1974年编制，通过男性特质、女性特质和中性特质三个维度评估个体的性别角色认同与心理双性化程度，打破传统二元性别观念，揭示性别角色对心理健康与社会适应的影响，是性别心理学研究领域的经典标准化工具。", []ReferenceItem{{Authors: []string{"Bem, S.L."}, Title: "The measurement of psychological androgyny", Journal: "Journal of Consulting and Clinical Psychology", Year: 1974, DOI: "10.1037/h0036215"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.BemSexRoleInventory,
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