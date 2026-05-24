package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func PhysicalAggression() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "我不喜欢受到限制或被命令做事。"},
		{Text: "我有时候会故意伤害别人的感情。"},
		{Text: "我常常感到无聊，需要一些刺激性的活动。"},
		{Text: "我觉得自己比别人更有权利和优越感。"},
		{Text: "我常常为了自己的利益而撒谎或欺骗别人。"},
		{Text: "我对别人的感受不太关心，也不愿意帮助他人。"},
		{Text: "我有时会因为自己的情绪而变得易怒和暴躁。"},
		{Text: "我常常追求冒险和刺激，不考虑后果。"},
		{Text: "我常常觉得自己比别人聪明。"},
		{Text: "你主要靠直觉和个人感情解决问题"},
		{Text: "别人对你的谈论不会使你难受"},
		{Text: "你喜欢谈不可能发生的事"},
		{Text: "你不善于修修补补"},
		{Text: "抽象地想象有助于解决问题"},
		{Text: "你喜欢事前将事情准备好"},
		{Text: "你不善于令别人相信"},
		{Text: "想个新主意对你来说不成问题"},
		{Text: "你几乎总是自由地作出决定"},
		{Text: "辩论中，你善于抓别人的弱点"},
	}

	// 定义选项
	options := []Option{
		{Tag: "A", Content: "否", Score: 0},
		{Tag: "B", Content: "是", Score: 1},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)

	// 保存测试类型信息
	SaveTestType(ctx, consts.PhysicalAggression, "躯体攻击性测试", "躯体攻击性量表源自Buss和Perry于1992年编制的攻击性问卷(AQ)，专门评估个体使用身体力量伤害他人的行为倾向与冲动控制能力，是攻击行为四维度模型（躯体攻击、言语攻击、愤怒、敌意）中的核心分量表，广泛应用于暴力行为预测与临床干预研究。", []ReferenceItem{{Authors: []string{"Buss, A.H.", "Perry, M."}, Title: "The aggression questionnaire", Journal: "Journal of Personality and Social Psychology", Year: 1992, DOI: "10.1037/0022-3514.63.3.452"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.PhysicalAggression).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.PhysicalAggression,
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
