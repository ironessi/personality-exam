package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// BDI 贝克抑郁自评量表初始化
func BDI() {
	ctx := context.Background()

	// 定义选项（BDI使用0-3分计分）
	options := []Option{
		{Tag: "A", Content: "无", Score: 0},
		{Tag: "B", Content: "轻度", Score: 1},
		{Tag: "C", Content: "中度", Score: 2},
		{Tag: "D", Content: "重度", Score: 3},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我不感到悲伤"},
		{Text: "我感到悲伤"},
		{Text: "我一直感到悲伤，不能自拔"},
		{Text: "我感到极度悲伤，甚至不能自制"},
		{Text: "我对未来并不悲观失望"},
		{Text: "我对未来感到悲观"},
		{Text: "我感到前景黯淡"},
		{Text: "我觉得将来毫无希望，无法改善"},
		{Text: "我没有失败的感觉"},
		{Text: "我觉得和大多数人相比我是失败的"},
		{Text: "回顾我的一生，我觉得那是一连串的失败"},
		{Text: "我觉得我是个彻底失败的人"},
		{Text: "我对事物像往常一样满意"},
		{Text: "我对事物不像往常一样满意"},
		{Text: "我不再对任何事情感到真正的满意"},
		{Text: "我对一切事物都不满意或讨厌"},
		{Text: "我没有特别感到内疚"},
		{Text: "相当多的时间我感到内疚"},
		{Text: "我大部分时间感到内疚"},
		{Text: "我一直感到内疚"},
		{Text: "我没有感到正在受惩罚"},
		{Text: "我感到我可能受惩罚"},
		{Text: "我预感会受惩罚"},
		{Text: "我感到我正在受惩罚"},
		{Text: "我对自己并不失望"},
		{Text: "我对自己感到失望"},
		{Text: "我讨厌自己"},
		{Text: "我憎恨自己"},
		{Text: "我没有责备自己把事情弄糟了"},
		{Text: "我以各种方式责备自己"},
		{Text: "我总是责备自己"},
		{Text: "我责备自己所有的事情都弄坏了"},
		{Text: "我没有任何自杀的想法"},
		{Text: "我有自杀的想法，但我不会去做"},
		{Text: "我想自杀"},
		{Text: "如果有机会我就会自杀"},
		{Text: "我并不比以往爱哭"},
		{Text: "我现在比以前爱哭"},
		{Text: "现在我经常哭"},
		{Text: "我以往能哭，但现在即使我想哭也哭不出来"},
		{Text: "我并不比以往容易生气"},
		{Text: "我现在比以前容易生气"},
		{Text: "我现在经常生气"},
		{Text: "以往能控制我的脾气，现在我动不动就生气，很难控制"},
		{Text: "我对他人的兴趣没有减少"},
		{Text: "我对他人的兴趣比以往减少了"},
		{Text: "我对他人丧失了大部分兴趣"},
		{Text: "我对他人现在毫无兴趣"},
		{Text: "我与以往一样能作决定"},
		{Text: "我现在作决定没有以前果断"},
		{Text: "我现在作决定比以前困难得多"},
		{Text: "我现在完全不能作决定"},
		{Text: "我觉得自己看上去和以前差不多"},
		{Text: "我担心我看上去老了，不吸引人了"},
		{Text: "我觉得我的外貌变得比以前难看了"},
		{Text: "我相信我看起来很丑陋"},
		{Text: "我能像以往一样工作"},
		{Text: "我做事比以前慢了"},
		{Text: "我做任何事都必须很努力地催促自己才行"},
		{Text: "我完全不能工作"},
		{Text: "我睡眠像以往一样好"},
		{Text: "我睡眠不如以往好"},
		{Text: "我比往常早醒1—2小时，难以再睡"},
		{Text: "我比往常早醒几个小时，不能再睡"},
		{Text: "我现在并不比以往感到容易疲劳"},
		{Text: "我现在比以往容易疲劳"},
		{Text: "我做任何事都容易疲劳"},
		{Text: "我太疲劳了，以致我不能做任何事情"},
		{Text: "我的食欲和以往一样好"},
		{Text: "我的食欲不如以往好"},
		{Text: "我的食欲差很多"},
		{Text: "我完全没有食欲"},
		{Text: "最近我的体重并无很大减轻"},
		{Text: "我的体重下降了2.27公斤以上"},
		{Text: "我的体重下降了4.54公斤以上"},
		{Text: "我的体重下降了6.8公斤以上"},
		{Text: "我对健康状况并不比以往更担心"},
		{Text: "我担心身体上的问题，如疼痛、胃不适或便秘"},
		{Text: "我很担心身体问题，想别的事情很难"},
		{Text: "我对身体问题如此担忧，以致不能想其他任何事情"},
		{Text: "我没有发现自己对性的兴趣最近有什么变化"},
		{Text: "我对性的兴趣比以往降低了"},
		{Text: "我现在对性的兴趣大大下降"},
		{Text: "我对性完全失去了兴趣"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.BDI), "贝克抑郁自评量表(BDI)", "贝克抑郁量表(BDI)由Beck等于1961年编制，包含21个项目，涵盖悲伤、悲观、失败感、失去愉悦感、罪恶感、自我厌恶等认知、情感与躯体症状维度，是全球临床心理学中最广泛使用的抑郁自评工具，适用于抑郁症状筛查、严重程度评估与治疗效果追踪。", []ReferenceItem{{Authors: []string{"Beck, A.T.", "Ward, C.H.", "Mendelson, M.", "Mock, J.", "Erbaugh, J."}, Title: "An inventory for measuring depression", Journal: "Archives of General Psychiatry", Year: 1961, DOI: "10.1001/archpsyc.1961.01710120031004"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.BDI,
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