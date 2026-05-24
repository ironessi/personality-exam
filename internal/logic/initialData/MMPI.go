package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// MMPI 多相人格测试初始化
func MMPI() {
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
		// 疑病量表
		{Text: "我经常感到身体不适"},
		{Text: "我担心自己的健康状况"},
		{Text: "我经常感到疲劳"},
		{Text: "我有很多身体上的抱怨"},
		{Text: "我对自己的健康非常关注"},
		{Text: "我经常感到头痛"},
		{Text: "我经常感到胃痛"},
		{Text: "我认为自己的健康状况不如他人"},
		{Text: "我经常感到头晕"},
		{Text: "我对自己的身体状况感到焦虑"},
		
		// 抑郁量表
		{Text: "我感到情绪低落"},
		{Text: "我对生活失去了兴趣"},
		{Text: "我感到绝望"},
		{Text: "我感到无助"},
		{Text: "我感到生活没有意义"},
		{Text: "我感到自己没有价值"},
		{Text: "我经常感到悲伤"},
		{Text: "我对未来感到悲观"},
		{Text: "我感到孤独"},
		{Text: "我很难感到快乐"},
		
		// 癔症量表
		{Text: "我在压力下会感到身体不适"},
		{Text: "我有时会失去记忆"},
		{Text: "我对疼痛非常敏感"},
		{Text: "我在情绪激动时会感到身体不适"},
		{Text: "我有时会感到虚弱无力"},
		{Text: "我对批评非常敏感"},
		{Text: "我在紧张时会感到身体不适"},
		{Text: "我有时会感到呼吸困难"},
		{Text: "我对压力的反应很强烈"},
		{Text: "我有时会感到心跳加速"},
		
		// 精神病态量表
		{Text: "我曾经违反过规则"},
		{Text: "我有时会撒谎"},
		{Text: "我对他人的感受不太敏感"},
		{Text: "我曾经有过违法行为"},
		{Text: "我有时会冲动行事"},
		{Text: "我不太关心社会规范"},
		{Text: "我曾经与人发生过冲突"},
		{Text: "我有时会感到愤怒"},
		{Text: "我不太在意他人的看法"},
		{Text: "我曾经有过不负责任的行为"},
		
		// 男性化-女性化量表
		{Text: "我对传统男性/女性角色感到舒适"},
		{Text: "我有符合自己性别的兴趣爱好"},
		{Text: "我认同自己的性别身份"},
		{Text: "我有符合自己性别的行为方式"},
		{Text: "我对异性有正常的兴趣"},
		{Text: "我有符合自己性别的情感表达"},
		{Text: "我对自己的性别特征感到满意"},
		{Text: "我有符合自己性别的职业兴趣"},
		{Text: "我对性别角色的期望感到舒适"},
		{Text: "我有符合自己性别的社交行为"},
		
		// 妄想狂量表
		{Text: "我认为有人在监视我"},
		{Text: "我认为有人在议论我"},
		{Text: "我认为有人想要伤害我"},
		{Text: "我对他人的动机感到怀疑"},
		{Text: "我认为自己被误解了"},
		{Text: "我认为有人在背后说我坏话"},
		{Text: "我对他人的信任度较低"},
		{Text: "我认为自己受到了不公正的待遇"},
		{Text: "我对他人的行为感到怀疑"},
		{Text: "我认为自己的想法被他人知道了"},
		
		// 精神衰弱量表
		{Text: "我经常感到焦虑"},
		{Text: "我经常感到紧张"},
		{Text: "我经常担心未来"},
		{Text: "我很难集中注意力"},
		{Text: "我经常感到不安"},
		{Text: "我经常感到害怕"},
		{Text: "我很难做出决定"},
		{Text: "我经常感到内疚"},
		{Text: "我经常感到自责"},
		{Text: "我经常感到担忧"},
		
		// 精神分裂症量表
		{Text: "我有时会有奇怪的想法"},
		{Text: "我有时会听到别人听不到的声音"},
		{Text: "我有时会看到别人看不到的东西"},
		{Text: "我有时会感到与现实脱节"},
		{Text: "我有时会有不寻常的体验"},
		{Text: "我有时会感到困惑"},
		{Text: "我有时会有奇怪的感觉"},
		{Text: "我有时会有不寻常的想法"},
		{Text: "我有时会感到自己与他人不同"},
		{Text: "我有时会有奇怪的行为"},
		
		// 轻躁狂量表
		{Text: "我有时会感到异常兴奋"},
		{Text: "我有时会有过多的能量"},
		{Text: "我有时会说话过多"},
		{Text: "我有时会有过多的想法"},
		{Text: "我有时会感到自己很重要"},
		{Text: "我有时会做出冲动的决定"},
		{Text: "我有时会感到自己充满创造力"},
		{Text: "我有时会睡眠减少但仍感到精力充沛"},
		{Text: "我有时会有夸大的想法"},
		{Text: "我有时会感到自己无所不能"},
		
		// 社会内向性量表
		{Text: "我喜欢独处"},
		{Text: "我不太喜欢社交活动"},
		{Text: "我对与人交往感到紧张"},
		{Text: "我宁愿独自工作"},
		{Text: "我不太喜欢与人交谈"},
		{Text: "我对社交场合感到不舒适"},
		{Text: "我很少主动与人联系"},
		{Text: "我对陌生人感到紧张"},
		{Text: "我不太喜欢参加聚会"},
		{Text: "我宁愿独自度过空闲时间"},
		
		// 补充问题
		{Text: "我对自己的生活感到满意"},
		{Text: "我有良好的人际关系"},
		{Text: "我能够应对生活中的挑战"},
		{Text: "我有明确的人生目标"},
		{Text: "我感到自己被他人接受"},
		{Text: "我有良好的自我形象"},
		{Text: "我能够表达自己的情感"},
		{Text: "我有良好的判断力"},
		{Text: "我能够适应变化"},
		{Text: "我感到生活有意义"},
		{Text: "我有良好的情绪调节能力"},
		{Text: "我能够与他人建立亲密关系"},
		{Text: "我有良好的问题解决能力"},
		{Text: "我感到自己有价值"},
		{Text: "我能够应对压力"},
		{Text: "我有良好的社交技能"},
		{Text: "我能够信任他人"},
		{Text: "我有良好的自我控制能力"},
		{Text: "我能够接受自己的缺点"},
		{Text: "我有良好的生活习惯"},
		{Text: "我能够享受生活"},
		{Text: "我有良好的健康习惯"},
		{Text: "我能够与他人合作"},
		{Text: "我有良好的时间管理能力"},
		{Text: "我能够设定并实现目标"},
		{Text: "我有良好的沟通能力"},
		{Text: "我能够处理冲突"},
		{Text: "我有良好的适应能力"},
		{Text: "我能够保持积极的态度"},
		{Text: "我有良好的自我意识"},
		{Text: "我能够与他人建立信任关系"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.MMPI), "多相人格测试", "明尼苏达多项人格测验(MMPI)由Hathaway和McKinley于1943年编制，是世界上使用最广泛的标准化人格测验之一，通过多个临床量表（抑郁、癔症、偏执等）和效度量表综合评估个体的人格特征与心理病理状态，广泛应用于临床心理诊断、司法鉴定与人员选拔领域。", []ReferenceItem{{Authors: []string{"Hathaway, S.R.", "McKinley, J.C."}, Title: "The Minnesota Multiphasic Personality Inventory", Journal: "University of Minnesota Press", Year: 1943}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.MMPI,
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