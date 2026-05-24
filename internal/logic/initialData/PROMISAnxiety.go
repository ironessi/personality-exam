package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func PROMISAnxiety() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "我担心工作或学习表现不好。"},
		{Text: "我感到喉咙里有梗塞感。"},
		{Text: "我容易感到疲劳，精力不足。"},
		{Text: "我担心家人或朋友的安全。"},
		{Text: "我出现胃痛、消化不良或腹泻等肠胃不适。"},
		{Text: "我做决定时变得非常困难。"},
		{Text: "我因为焦虑而避免去某些地方或场合。"},
		{Text: "我感到一阵阵的发热或发冷。"},
		{Text: "我对自己失去信心。"},
		{Text: "我担心未来的事情。"},
		{Text: "我感到无助，觉得事情难以应付。"},
		{Text: "我因为烦躁而对他人的容忍度降低。"},
		{Text: "我感到不真实，或者与周围环境有隔阂。"},
		{Text: "我担心自己会突然崩溃。"},
		{Text: "我坐立不安，需要来回走动。"},
		{Text: "我因为思虑过重而感到头脑混乱。"},
		{Text: "我担心自己得不到他人的认可。"},
		{Text: "总的来说，我感到持续的焦虑和担忧。"},
		{Text: "我发现自己很难安静地坐着。"},
		{Text: "我担心一些实际上不会发生的事。"},
		{Text: "我时常无缘无故地感到害怕。"},
		{Text: "我的肌肉紧张，感觉僵硬和酸痛。"},
		{Text: "我很难放松下来。"},
		{Text: "我因为担心而难以入睡，或睡眠很浅。"},
		{Text: "我担心自己的健康状况。"},
		{Text: "我感到烦躁，容易因为小事受惊吓。"},
		{Text: "我的心跳得很快，或者感觉心律不齐。"},
		{Text: "我经常感到紧张和心神不定。"},
		{Text: "我担心会在大庭广众下出丑或失控。"},
		{Text: "我呼吸变得急促或感到呼吸困难。"},
		{Text: "我的手心常常出汗或者手脚冰凉。"},
		{Text: "我感到头晕或头昏。"},
		{Text: "我担心经济状况或财务状况。"},
		{Text: "我感到颤抖（例如手抖）。"},
		{Text: "我因为焦虑而难以集中注意力。"},
	}

	// 定义选项
	options := []Option{
		{Tag: "A", Content: "不符合", Score: 0},
		{Tag: "B", Content: "有点符合", Score: 1},
		{Tag: "C", Content: "比较符合", Score: 2},
		{Tag: "D", Content: "非常符合", Score: 3},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.PROMISAnxiety), "忧虑", "PROMIS焦虑量表由美国国立卫生研究院于2004年作为患者报告结局测量信息系统的核心模块开发，基于项目反应理论构建，精准评估个体在过去七天内的恐惧感、忧虑程度、过度警觉及躯体焦虑症状，适用于临床研究、疗效追踪与患者报告结局的标准化测量。", []ReferenceItem{{Authors: []string{"Cella, D.", "Yount, S.", "Rothrock, N."}, Title: "The Patient-Reported Outcomes Measurement Information System (PROMIS)", Journal: "Medical Care", Year: 2007, DOI: "10.1097/01.mlr.0000258615.42478.55"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.PROMISAnxiety).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.PROMISAnxiety,
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
			log.Printf("插入选项数据失败: %v\n", err)
			return
		}
	}
}
