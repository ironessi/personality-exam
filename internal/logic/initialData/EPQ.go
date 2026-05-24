package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// EPQ EPQ艾森克人格测试初始化
func EPQ() {
	ctx := context.Background()

	// 定义选项（是/否计分）
	options := []Option{
		{Tag: "A", Content: "是", Score: 1},
		{Tag: "B", Content: "否", Score: 0},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "你是否是一个健谈的人？"},
		{Text: "你是否喜欢参加大型社交活动？"},
		{Text: "你是否容易感到紧张或焦虑？"},
		{Text: "你是否经常担心可能发生的事情？"},
		{Text: "你是否喜欢独处？"},
		{Text: "你是否容易感到沮丧或忧郁？"},
		{Text: "你是否善于结交新朋友？"},
		{Text: "你是否经常感到不安？"},
		{Text: "你是否喜欢冒险？"},
		{Text: "你是否容易激动？"},
		{Text: "你是否经常感到疲倦？"},
		{Text: "你是否喜欢尝试新事物？"},
		{Text: "你是否容易感到愤怒？"},
		{Text: "你是否喜欢与他人合作？"},
		{Text: "你是否经常感到孤独？"},
		{Text: "你是否善于表达自己的情感？"},
		{Text: "你是否经常感到担忧？"},
		{Text: "你是否喜欢热闹的聚会？"},
		{Text: "你是否容易感到压力？"},
		{Text: "你是否喜欢独立工作？"},
		{Text: "你是否经常感到烦躁？"},
		{Text: "你是否善于说服他人？"},
		{Text: "你是否经常感到害怕？"},
		{Text: "你是否喜欢旅行？"},
		{Text: "你是否容易感到悲伤？"},
		{Text: "你是否喜欢与人交谈？"},
		{Text: "你是否经常感到焦虑？"},
		{Text: "你是否喜欢参加体育活动？"},
		{Text: "你是否容易感到不安？"},
		{Text: "你是否喜欢阅读？"},
		{Text: "你是否经常感到情绪低落？"},
		{Text: "你是否喜欢与他人分享？"},
		{Text: "你是否经常感到紧张？"},
		{Text: "你是否喜欢参加集体活动？"},
		{Text: "你是否容易感到疲惫？"},
		{Text: "你是否喜欢听音乐？"},
		{Text: "你是否经常感到忧虑？"},
		{Text: "你是否喜欢看电影？"},
		{Text: "你是否容易感到烦躁不安？"},
		{Text: "你是否喜欢购物？"},
		{Text: "你是否经常感到压力很大？"},
		{Text: "你是否喜欢烹饪？"},
		{Text: "你是否容易感到沮丧？"},
		{Text: "你是否喜欢园艺？"},
		{Text: "你是否经常感到焦虑不安？"},
		{Text: "你是否喜欢绘画？"},
		{Text: "你是否容易感到紧张不安？"},
		{Text: "你是否喜欢写作？"},
		{Text: "你是否经常感到情绪不稳定？"},
		{Text: "你是否喜欢摄影？"},
		{Text: "你是否容易感到愤怒或生气？"},
		{Text: "你是否喜欢钓鱼？"},
		{Text: "你是否经常感到害怕或恐惧？"},
		{Text: "你是否喜欢爬山？"},
		{Text: "你是否容易感到悲伤或忧郁？"},
		{Text: "你是否喜欢游泳？"},
		{Text: "你是否经常感到焦虑或紧张？"},
		{Text: "你是否喜欢跑步？"},
		{Text: "你是否容易感到疲惫或劳累？"},
		{Text: "你是否喜欢骑自行车？"},
		{Text: "你是否经常感到烦躁或易怒？"},
		{Text: "你是否喜欢瑜伽？"},
		{Text: "你是否经常感到情绪低落或沮丧？"},
		{Text: "你是否喜欢冥想？"},
		{Text: "你是否容易感到压力或紧张？"},
		{Text: "你是否喜欢跳舞？"},
		{Text: "你是否经常感到不安或焦虑？"},
		{Text: "你是否喜欢唱歌？"},
		{Text: "你是否容易感到愤怒或生气？"},
		{Text: "你是否喜欢演奏乐器？"},
		{Text: "你是否经常感到害怕或恐惧？"},
		{Text: "你是否喜欢戏剧？"},
		{Text: "你是否容易感到悲伤或忧郁？"},
		{Text: "你是否喜欢歌剧？"},
		{Text: "你是否经常感到焦虑或紧张？"},
		{Text: "你是否喜欢芭蕾？"},
		{Text: "你是否容易感到疲惫或劳累？"},
		{Text: "你是否喜欢杂技？"},
		{Text: "你是否经常感到烦躁或易怒？"},
		{Text: "你是否喜欢魔术？"},
		{Text: "你是否经常感到情绪低落或沮丧？"},
		{Text: "你是否喜欢马戏？"},
		{Text: "你是否容易感到压力或紧张？"},
		{Text: "你是否喜欢木偶戏？"},
		{Text: "你是否经常感到不安或焦虑？"},
		{Text: "你是否喜欢皮影戏？"},
		{Text: "你是否容易感到愤怒或生气？"},
		{Text: "你是否喜欢漫画？"},
		{Text: "你是否经常感到害怕或恐惧？"},
		{Text: "你是否喜欢动画？"},
		{Text: "你是否容易感到悲伤或忧郁？"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.EPQ), "EPQ艾森克人格测试-国际标准版", "艾森克人格问卷(EPQ)由英国心理学家Hans Eysenck于1975年编制，基于三维人格理论，通过神经质(N)、外向性(E)、精神质(P)三个核心维度及谎言量表(L)全面评估个体人格特征，具有良好的跨文化信效度，是人格心理学研究与临床评估领域的经典标准化工具。", []ReferenceItem{{Authors: []string{"Eysenck, H.J.", "Eysenck, S.B.G."}, Title: "Manual of the Eysenck Personality Questionnaire", Journal: "Hodder and Stoughton", Year: 1975}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.EPQ,
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