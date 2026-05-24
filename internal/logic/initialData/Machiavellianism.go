package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// Machiavellianism 马基雅维利人格测试初始化
func Machiavellianism() {
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
		{Text: "为了达到目的，我可以不择手段"},
		{Text: "我认为诚实是最好的策略"},
		{Text: "我喜欢控制他人"},
		{Text: "我认为人大多是自私的"},
		{Text: "我善于操纵他人"},
		{Text: "我认为道德是相对的"},
		{Text: "我喜欢权力"},
		{Text: "我认为为了成功可以牺牲道德"},
		{Text: "我善于说谎"},
		{Text: "我认为人应该为了自己的利益而奋斗"},
		{Text: "我喜欢利用他人"},
		{Text: "我认为忠诚是一种弱点"},
		{Text: "我善于说服他人"},
		{Text: "我认为人大多是不可信的"},
		{Text: "我喜欢在背后操纵事情"},
		{Text: "我认为为了达到目的可以欺骗他人"},
		{Text: "我善于隐藏自己的真实意图"},
		{Text: "我认为人应该为了自己的目标而努力"},
		{Text: "我喜欢影响他人的决策"},
		{Text: "我认为道德规范是为了约束他人而存在的"},
		{Text: "我善于利用他人的弱点"},
		{Text: "我认为诚实并不总是最好的策略"},
		{Text: "我喜欢控制局势"},
		{Text: "我认为人大多是虚伪的"},
		{Text: "我善于为自己的行为找借口"},
		{Text: "我认为为了成功可以打破规则"},
		{Text: "我喜欢成为领导者"},
		{Text: "我认为人应该为了自己的利益而不择手段"},
		{Text: "我善于利用他人的信任"},
		{Text: "我认为道德是用来约束他人的"},
		{Text: "我喜欢操纵他人的情绪"},
		{Text: "我认为人大多是利己的"},
		{Text: "我善于伪装自己"},
		{Text: "我认为为了达到目的可以牺牲他人的利益"},
		{Text: "我喜欢掌控一切"},
		{Text: "我认为忠诚是一种愚蠢的表现"},
		{Text: "我善于利用他人的资源"},
		{Text: "我认为人应该为了自己的目标而不择手段"},
		{Text: "我喜欢在他人面前表现得自信"},
		{Text: "我认为道德是相对的，没有绝对的对错"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.Machiavellianism), "马基雅维利人格测试", "马基雅维利主义量表(Mach-IV)由Christie和Geis于1970年编制，评估个体操纵他人的倾向、对道德的冷漠态度及对权力的渴望程度，是人格心理学中暗黑三联征（马基雅维利主义、自恋、精神病态）研究的核心测量工具，广泛应用于组织行为学、政治心理学与人际关系研究。", []ReferenceItem{{Authors: []string{"Christie, R.", "Geis, F.L."}, Title: "Studies in Machiavellianism", Journal: "Academic Press", Year: 1970}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.Machiavellianism,
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