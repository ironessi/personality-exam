package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// EPPS 爱德华个人偏好测验初始化
func EPPS() {
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
		// 成就
		{Text: "我喜欢完成具有挑战性的任务"},
		{Text: "我追求卓越的表现"},
		{Text: "我为自己的成就感到自豪"},
		{Text: "我喜欢设定并实现目标"},
		
		// 顺从
		{Text: "我遵守规则和权威"},
		{Text: "我尊重传统价值观"},
		{Text: "我愿意服从他人的指导"},
		{Text: "我重视社会规范"},
		
		// 秩序
		{Text: "我喜欢整洁有序的环境"},
		{Text: "我按照计划和时间表行事"},
		{Text: "我重视组织和结构"},
		{Text: "我喜欢把事情安排得井井有条"},
		
		// 表现
		{Text: "我喜欢成为众人关注的焦点"},
		{Text: "我乐于在公众场合表现自己"},
		{Text: "我喜欢展示自己的才华"},
		{Text: "我享受被他人赞赏的感觉"},
		
		// 自主
		{Text: "我喜欢独立做决定"},
		{Text: "我重视个人自由"},
		{Text: "我不喜欢被他人控制"},
		{Text: "我喜欢按照自己的方式做事"},
		
		// 亲和
		{Text: "我喜欢与他人建立亲密关系"},
		{Text: "我重视友谊"},
		{Text: "我乐于帮助他人"},
		{Text: "我喜欢与他人合作"},
		
		// 省察
		{Text: "我经常反思自己的行为"},
		{Text: "我关注自己的内心感受"},
		{Text: "我喜欢分析自己的动机"},
		{Text: "我重视自我了解"},
		
		// 求助
		{Text: "我在需要时会寻求他人的帮助"},
		{Text: "我喜欢被他人关心和照顾"},
		{Text: "我重视他人的支持"},
		{Text: "我不喜欢独自面对困难"},
		
		// 支配
		{Text: "我喜欢领导和影响他人"},
		{Text: "我重视自己的影响力"},
		{Text: "我乐于做决定并承担责任"},
		{Text: "我喜欢掌控局面"},
		
		// 谦卑
		{Text: "我愿意承认自己的错误"},
		{Text: "我重视他人的意见"},
		{Text: "我不喜欢炫耀自己的成就"},
		{Text: "我乐于接受批评"},
		
		// 慈善
		{Text: "我喜欢帮助有需要的人"},
		{Text: "我重视他人的福祉"},
		{Text: "我乐于分享自己的资源"},
		{Text: "我关心弱势群体"},
		
		// 变通
		{Text: "我喜欢尝试新的事物"},
		{Text: "我适应变化的能力强"},
		{Text: "我重视灵活性"},
		{Text: "我不喜欢一成不变的生活"},
		
		// 坚毅
		{Text: "我能够坚持完成困难的任务"},
		{Text: "我在面对挫折时不轻易放弃"},
		{Text: "我重视毅力和恒心"},
		{Text: "我能够长时间专注于一项工作"},
		
		// 异性恋
		{Text: "我对异性有浪漫的兴趣"},
		{Text: "我重视与异性的关系"},
		{Text: "我希望与异性建立亲密关系"},
		{Text: "我对异性的吸引力敏感"},
		
		// 攻击
		{Text: "我在面对不公时会表达愤怒"},
		{Text: "我喜欢竞争和挑战"},
		{Text: "我不害怕与他人发生冲突"},
		{Text: "我重视自己的权益"},
		
		// 第二轮问题
		// 成就
		{Text: "我喜欢超越自己的极限"},
		{Text: "我追求个人成长和进步"},
		{Text: "我为实现目标而努力工作"},
		{Text: "我喜欢获得成功的感觉"},
		
		// 顺从
		{Text: "我尊重长辈和权威人士"},
		{Text: "我愿意遵守承诺"},
		{Text: "我重视责任和义务"},
		{Text: "我不喜欢违反规则"},
		
		// 秩序
		{Text: "我喜欢制定计划和安排"},
		{Text: "我重视准确性和 precision"},
		{Text: "我喜欢按照既定的程序做事"},
		{Text: "我不喜欢混乱和无序"},
		
		// 表现
		{Text: "我喜欢在社交场合中表现活跃"},
		{Text: "我乐于表达自己的想法和感受"},
		{Text: "我喜欢参与表演和展示活动"},
		{Text: "我不害怕被他人评价"},
		
		// 自主
		{Text: "我喜欢独立思考"},
		{Text: "我重视个人隐私"},
		{Text: "我不喜欢被他人打扰"},
		{Text: "我喜欢有自己的空间和时间"},
		
		// 亲和
		{Text: "我喜欢与他人交流和分享"},
		{Text: "我重视团队精神"},
		{Text: "我乐于参与集体活动"},
		{Text: "我不喜欢孤独"},
		
		// 省察
		{Text: "我喜欢分析问题和思考"},
		{Text: "我重视理性和逻辑"},
		{Text: "我乐于探索自己的内心世界"},
		{Text: "我不喜欢盲目跟从他人"},
		
		// 求助
		{Text: "我喜欢与他人合作解决问题"},
		{Text: "我重视团队的支持"},
		{Text: "我乐于接受他人的建议"},
		{Text: "我不喜欢独自承担压力"},
		
		// 支配
		{Text: "我喜欢组织和管理事务"},
		{Text: "我重视自己的领导力"},
		{Text: "我乐于指导和帮助他人成长"},
		{Text: "我不喜欢被他人领导"},
		
		// 谦卑
		{Text: "我喜欢学习他人的优点"},
		{Text: "我重视谦虚和低调"},
		{Text: "我乐于赞美他人的成就"},
		{Text: "我不喜欢自夸和炫耀"},
		
		// 慈善
		{Text: "我喜欢参与志愿服务"},
		{Text: "我重视同理心和 compassion"},
		{Text: "我乐于关心和照顾他人"},
		{Text: "我不喜欢看到他人受苦"},
		
		// 变通
		{Text: "我喜欢探索新的领域和机会"},
		{Text: "我重视创新和创造力"},
		{Text: "我乐于接受新的挑战"},
		{Text: "我不喜欢墨守成规"},
		
		// 坚毅
		{Text: "我喜欢克服困难和障碍"},
		{Text: "我重视勇气和决心"},
		{Text: "我乐于面对挑战"},
		{Text: "我不喜欢轻易妥协"},
		
		// 异性恋
		{Text: "我喜欢与异性交往和互动"},
		{Text: "我重视爱情和浪漫关系"},
		{Text: "我乐于与异性建立深度连接"},
		{Text: "我对异性的情感需求敏感"},
		
		// 攻击
		{Text: "我喜欢为自己的信仰而辩护"},
		{Text: "我重视正义和公平"},
		{Text: "我乐于挑战不合理的事情"},
		{Text: "我不喜欢被他人欺负或利用"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.EPPS), "爱德华个人偏好测验", "爱德华个人偏好测验(EPPS)由Edwards于1953年编制，基于Murray的心理需要理论，通过225对强迫选择题评估个体在成就、顺从、秩序、表现、自主等15种主要心理需求上的相对强度，有效控制社会期望偏差，揭示个体内在动机结构与人格倾向，广泛应用于职业咨询与人格研究。", []ReferenceItem{{Authors: []string{"Edwards, A.L."}, Title: "Edwards Personal Preference Schedule", Journal: "Psychological Corporation", Year: 1953}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.EPPS,
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