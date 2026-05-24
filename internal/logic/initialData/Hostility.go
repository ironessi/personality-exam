package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func Hostility() {
	ctx := context.Background()

	var questions = []Question{
		{
			Text: "邻居家的宠物经常在你家门口大小便，你会：",
			Options: []Option{
				{Tag: "A", Content: "直接找邻居理论，要求邻居管好宠物，否则就采取措施", Score: 3},
				{Tag: "B", Content: "自己清理干净，然后贴个纸条提醒邻居注意", Score: 1},
				{Tag: "C", Content: "买些驱宠物的物品放在家门口，让邻居的宠物不敢靠近", Score: 2},
				{Tag: "D", Content: "联合其他邻居一起向物业反映，让物业处理", Score: 0},
			},
		},
		{
			Text: "你在马路上被一辆车溅了一身水，司机没有停车道歉就开走了，你会：",
			Options: []Option{
				{Tag: "A", Content: "记住车牌号码，向交警部门举报司机的不文明驾驶行为", Score: 1},
				{Tag: "B", Content: "对着车开走的方向大骂几句，发泄情绪", Score: 3},
				{Tag: "C", Content: "默默自认倒霉，回家换衣服", Score: 0},
				{Tag: "D", Content: "下次遇到类似情况，也故意开车溅别人一身水", Score: 2},
			},
		},
		{
			Text: "你的合作伙伴在合作项目中偷懒，导致进度延迟，你会：",
			Options: []Option{
				{Tag: "A", Content: "严厉指责合作伙伴，要求对方立刻改正，否则就终止合作", Score: 3},
				{Tag: "B", Content: "减少自己在项目中的投入，以同样的方式 \"报复\" 对方", Score: 2},
				{Tag: "C", Content: "与合作伙伴坐下来，分析当前情况，明确各自的责任和任务", Score: 0},
				{Tag: "D", Content: "向其他合作方反映，让大家一起给这个合作伙伴压力", Score: 1},
			},
		},
		{
			Text: "合租室友经常带朋友回家，很吵闹，影响你休息，你提醒过几次都没用，你会：",
			Options: []Option{
				{Tag: "A", Content: "直接和室友吵架，要求室友以后不准带朋友回家", Score: 3},
				{Tag: "B", Content: "向房东反映情况，让房东出面解决", Score: 1},
				{Tag: "C", Content: "在室友带朋友回家时，故意制造噪音干扰他们", Score: 2},
				{Tag: "D", Content: "再次和室友心平气和地沟通，提出解决方案，比如规定带朋友回家的时间", Score: 0},
			},
		},
		{
			Text: "你参加比赛，结果被评委不公平打分，你认为自己的表现应该得到更高分数，你会：",
			Options: []Option{
				{Tag: "A", Content: "当场和评委理论，质疑评委的公正性", Score: 3},
				{Tag: "B", Content: "向比赛主办方投诉，要求重新评判", Score: 1},
				{Tag: "C", Content: "在社交媒体上吐槽评委，引发舆论关注", Score: 2},
				{Tag: "D", Content: "接受结果，但心里默默记下，以后有机会证明自己", Score: 0},
			},
		},
		{
			Text: "你发现自己的作品被他人抄袭，你会：",
			Options: []Option{
				{Tag: "A", Content: "直接联系抄袭者，要求对方立刻删除抄袭内容，并公开道歉", Score: 3},
				{Tag: "B", Content: "在社交平台曝光抄袭者的行为，让大家一起谴责", Score: 2},
				{Tag: "C", Content: "咨询法律专家，准备起诉抄袭者", Score: 1},
				{Tag: "D", Content: "先和抄袭者沟通，希望对方能主动承认错误并改正", Score: 0},
			},
		},
		{
			Text: "你在游戏中被队友坑了，导致游戏失败，你会：",
			Options: []Option{
				{Tag: "A", Content: "在游戏中疯狂辱骂队友，发泄不满", Score: 3},
				{Tag: "B", Content: "默默拉黑队友，以后不再和他一起玩游戏", Score: 1},
				{Tag: "C", Content: "向游戏官方举报队友的不良行为", Score: 2},
				{Tag: "D", Content: "提醒队友以后注意游戏操作，这次就算了", Score: 0},
			},
		},
		{
			Text: "亲戚借了你一笔钱，到期了还不还，你催了几次都无果，你会：",
			Options: []Option{
				{Tag: "A", Content: "去亲戚家大闹一场，让亲戚在家人面前难堪，逼迫还钱", Score: 3},
				{Tag: "B", Content: "向其他亲戚诉说，让大家给这个亲戚施加压力", Score: 2},
				{Tag: "C", Content: "找律师咨询，准备通过法律途径解决", Score: 1},
				{Tag: "D", Content: "再和亲戚好好沟通一次，说明自己的难处，希望对方还钱", Score: 0},
			},
		},
		{
			Text: "你在商场购物时，被导购员怀疑没付款就想离开，你会：",
			Options: []Option{
				{Tag: "A", Content: "非常生气，要求导购员公开道歉，否则就投诉", Score: 3},
				{Tag: "B", Content: "配合导购员检查，证明自己的清白，然后默默离开，心里不爽", Score: 1},
				{Tag: "C", Content: "向商场负责人反映导购员的无礼行为，要求商场给出说法", Score: 2},
				{Tag: "D", Content: "觉得导购员可能工作失误，解释清楚后就不再追究", Score: 0},
			},
		},
		{
			Text: "你发现同事在背后说你坏话，你会：",
			Options: []Option{
				{Tag: "A", Content: "当面质问同事，让同事给你一个说法", Score: 3},
				{Tag: "B", Content: "以其人之道还治其人之身，也在背后说同事坏话", Score: 2},
				{Tag: "C", Content: "找机会在大家面前展现自己的优点，让大家看到同事说的不是事实", Score: 1},
				{Tag: "D", Content: "尝试和同事沟通，化解矛盾，询问对方为什么这么做", Score: 0},
			},
		},
		{
			Text: "当你在排队时，有人插队到你前面，你的第一反应是：",
			Options: []Option{
				{Tag: "A", Content: "立刻出声制止，要求对方到后面排队", Score: 3},
				{Tag: "B", Content: "心里很生气，但选择默默忍受，事后抱怨", Score: 1},
				{Tag: "C", Content: "瞪对方一眼，小声嘀咕几句", Score: 2},
				{Tag: "D", Content: "寻找工作人员，让他们来处理插队问题", Score: 0},
			},
		},
		{
			Text: "你精心准备的方案被领导批评得一无是处，你觉得很委屈，因为你觉得领导误解了你，这时你会：",
			Options: []Option{
				{Tag: "A", Content: "当场和领导理论，指出领导的错误理解", Score: 3},
				{Tag: "B", Content: "默默流泪，心里怨恨领导，工作积极性大减", Score: 1},
				{Tag: "C", Content: "私下向其他同事抱怨领导的不公，希望得到支持", Score: 2},
				{Tag: "D", Content: "冷静下来，找个合适的时间和领导再次沟通，解释方案的亮点", Score: 0},
			},
		},
		{
			Text: "你在餐厅用餐，服务员上菜很慢，还态度不好，你会：",
			Options: []Option{
				{Tag: "A", Content: "向餐厅经理投诉，要求对服务员进行处罚", Score: 2},
				{Tag: "B", Content: "给餐厅差评，在网上吐槽这家餐厅的服务", Score: 1},
				{Tag: "C", Content: "直接对服务员发火，要求快点上菜并改善态度", Score: 3},
				{Tag: "D", Content: "心平气和地提醒服务员加快速度，注意态度", Score: 0},
			},
		},
		{
			Text: "同学在班级聚会上开了一个让你很尴尬的玩笑，你会：",
			Options: []Option{
				{Tag: "A", Content: "当场翻脸，也开一个类似的玩笑回敬同学，让对方也难堪", Score: 3},
				{Tag: "B", Content: "尴尬地笑一笑，事后找机会单独和同学说这个玩笑不合适", Score: 1},
				{Tag: "C", Content: "记在心里，之后在同学聚会时找机会让对方出丑", Score: 2},
				{Tag: "D", Content: "理解同学可能没有恶意，一笑而过，不放在心上", Score: 0},
			},
		},
		{
			Text: "你的恋人与异性朋友走得很近，让你产生误会，你会：",
			Options: []Option{
				{Tag: "A", Content: "马上质问恋人，要求恋人立刻与异性朋友保持距离，否则就分手", Score: 3},
				{Tag: "B", Content: "旁敲侧击地提醒恋人，表达自己的感受，希望对方注意", Score: 1},
				{Tag: "C", Content: "偷偷观察恋人与异性朋友的互动，寻找证据，准备 \"秋后算账\"", Score: 2},
				{Tag: "D", Content: "选择信任恋人，找个合适的时机和恋人坦诚沟通这件事", Score: 0},
			},
		},
		{
			Text: "乘坐公共交通工具时，有人不小心踩了你一脚，还不道歉，你会：",
			Options: []Option{
				{Tag: "A", Content: "大声指责对方，让对方道歉", Score: 3},
				{Tag: "B", Content: "回踩对方一脚，以牙还牙", Score: 2},
				{Tag: "C", Content: "默默忍受，心里想着今天真倒霉", Score: 0},
				{Tag: "D", Content: "看着对方，说一句 \"你踩到我了\"，等待对方反应", Score: 1},
			},
		},
		{
			Text: "你借给邻居东西，对方很久都没归还，你去索要时，邻居态度还不好，这时你会：",
			Options: []Option{
				{Tag: "A", Content: "大发脾气，指责邻居不懂礼貌和感恩，要求立刻归还", Score: 3},
				{Tag: "B", Content: "向其他邻居吐槽这件事，让大家评评理", Score: 2},
				{Tag: "C", Content: "心里决定以后再也不借东西给对方，这次先把东西要回来", Score: 1},
				{Tag: "D", Content: "委婉地再次提醒邻居归还，不与对方起冲突", Score: 0},
			},
		},
		{
			Text: "假如你在网上看到有人恶意评论你的动态，你会：",
			Options: []Option{
				{Tag: "A", Content: "马上回怼，用同样难听的话攻击对方", Score: 3},
				{Tag: "B", Content: "选择举报对方的评论，让平台处理", Score: 2},
				{Tag: "C", Content: "不回复，但在自己后续动态中隐晦表达不满", Score: 1},
				{Tag: "D", Content: "直接忽略，觉得没必要和这种人计较", Score: 0},
			},
		},
		{
			Text: "同事在工作中抢了你的功劳，你得知后会：",
			Options: []Option{
				{Tag: "A", Content: "直接去找同事对峙，在众人面前揭露他的行为", Score: 3},
				{Tag: "B", Content: "向领导反映真实情况，有理有据地说明功劳归属", Score: 1},
				{Tag: "C", Content: "表面不动声色，私下寻找机会在工作中让对方难堪", Score: 2},
				{Tag: "D", Content: "认为这次吃亏就算了，以后注意保护自己的成果，不采取行动报复", Score: 0},
			},
		},
		{
			Text: "朋友不小心弄坏了你心爱的物品，你会：",
			Options: []Option{
				{Tag: "A", Content: "非常生气，要求对方立刻赔偿，态度强硬", Score: 3},
				{Tag: "B", Content: "虽然心疼，但觉得朋友不是故意的，安慰对方别在意", Score: 0},
				{Tag: "C", Content: "心里有些不悦，暗示朋友应该赔偿", Score: 2},
				{Tag: "D", Content: "先看看物品损坏程度，再决定是否让朋友赔偿，语气平和", Score: 1},
			},
		},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.Hostility), "敌意倾向测试", "敌意倾向测试基于Cook-Medley敌意量表，由Cook和Medley于1954年从MMPI中提取编制，评估个体对他人的愤世嫉俗态度、愤怒情绪及攻击性反应模式，研究表明高敌意得分与心血管疾病风险显著相关，是人格与健康心理学交叉研究的重要测量工具。", []ReferenceItem{{Authors: []string{"Cook, W.W.", "Medley, D.M."}, Title: "Proposed hostility and pharisaic-virtue scales for the MMPI", Journal: "Journal of Applied Psychology", Year: 1954, DOI: "10.1037/h0060667"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.Hostility).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.Hostility,
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

		// 将选项转换为JSON
		optionStr := optionsToJSON(question.Options)

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
