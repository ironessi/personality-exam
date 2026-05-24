package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func SCL90() {
	ctx := context.Background()

	questionContents := []string{
		"头痛", "无缘无故感到害怕", "反复思考无意义的事", "腰痛", "食欲不振",
		"对他人批评敏感", "听到他人听不见的声音", "感到有人监视自己", "担心衣饰不整或仪态失礼", "做事必须反复检查",
		"容易烦恼或激动", "头晕或昏倒", "害怕空旷场所", "常想哭泣", "难以做出决定",
		"感到思维被外界控制", "心跳突然加速", "认为有人企图伤害自己", "感到难以完成任务", "感到生活无意义",
		"与人交往时自卑", "胸痛", "坐立不安", "想摔坏物品", "害怕独自出门",
		"对未来绝望", "恶心或胃部不适", "脑中反复出现不必要的想法", "容易自责", "感到孤独",
		"情绪低沉", "突发性恐慌", "对脏污或病菌过分在意", "感到被他人轻视", "有超自然体验",
		"与异性相处不自在", "感到他人不理解自己", "难以集中注意力", "手脚发抖", "胸痛",
		"担心自己言行伤害他人", "肌肉酸痛", "感到被阴谋针对", "感到孤独", "出现无意义的仪式行为",
		"担心身体有大病", "害怕乘坐交通工具", "睡眠质量差", "呼吸困难", "害怕特定动物/场景",
		"控制不住想骂人", "做事兴趣下降", "身体发麻或刺痛", "比平时话少", "忘记性大",
		"感到呼吸困难", "喉咙有梗塞感", "有轻生念头", "感到内疚", "感到精力下降",
		"在人群中有孤立感", "感到他人能读取自己思想", "常有报复念头", "对事物过度担忧", "相信他人能操控自己行为",
		"感到身体衰弱无力", "手心潮湿或出汗", "怀疑伴侣不忠", "手脚沉重", "在人群中有窒息感",
		"认为他人窃取自己的想法", "感到他人对自己有敌意", "害怕高处或电梯", "面部潮红或发冷", "感到自己毫无价值",
		"需要频繁排尿", "独自发笑或自言自语", "对他人过分警惕", "认为自己是特殊人物", "害怕黑暗环境",
		"不愿透露真实想法", "感到周围事物不真实", "有怪异的身体感知", "夜间易醒或早醒", "担心自己会当众失控",
		"容易因小事羞耻", "对噪音敏感", "认为自己有特殊能力", "需要频繁排尿", "有怪异的身体感知",
		"对亲友发脾气", "认为自己是特殊人物", "害怕黑暗环境", "不愿透露真实想法", "感到大难临头",
	}

	// 定义选项
	options := []Option{
		{Tag: "A", Content: "没有", Score: 1},
		{Tag: "B", Content: "很轻", Score: 2},
		{Tag: "C", Content: "中度", Score: 3},
		{Tag: "D", Content: "偏重", Score: 4},
		{Tag: "E", Content: "严重", Score: 5},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)

	// 保存测试类型信息
	SaveTestType(ctx, consts.SCL90, "SCL90", "症状自评量表(SCL-90)由Derogatis于1973年编制，包含90个项目，涵盖躯体化、强迫症状、人际敏感、抑郁、焦虑、敌对、恐怖、偏执及精神病性9个维度，是目前国际上最广泛使用的心理健康自评工具之一，适用于临床筛查、疗效评估与流行病学研究。", []ReferenceItem{{Authors: []string{"Derogatis, L.R.", "Lipman, R.S.", "Covi, L."}, Title: "SCL-90: An outpatient psychiatric rating scale", Journal: "Psychopharmacology Bulletin", Year: 1973}}, len(questionContents))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.SCL90).Delete()
	if err != nil {
		log.Printf("删除问题失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, content := range questionContents {
		sequence := idx + 1

		// 插入问题
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.SCL90,
			Content:    content,
			Sequence:   sequence,
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

		// 插入选项 - 使用正确的JSON格式
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
