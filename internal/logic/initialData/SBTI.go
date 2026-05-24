package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// SBTI 傻大个性格测试初始化
func SBTI() {
	ctx := context.Background()

	var questions = []Question{
		{
			Text: "我是一只阴暗的老鼠，一只爬行的蟑螂，这辈子没谈过一场恋爱，胆怯又自卑，我的青春就是一场又一场的幻想，每一天幻想着我也能有一个爱我的人和我一起压马路，一起逛街，一起玩，现实却是花了父母的钱读了个学校，毕业之后找班上，寻找着理想，寻找着目标，每次看到大家能在网上开幸福的玩笑，我都想哭。我在羡慕生活的各种美好，却还在寻找自己的方向。每一次看到这种都是对我心灵的一次伤害，求求大家给我们这种在被窝里偷窥的人一点活路吧，我真的不想在白天把枕巾哭湿一大片",
			Options: []Option{
				{Content: "我哭了。。", Score: 1},
				{Content: "这是什么。。", Score: 2},
				{Content: "这不是我！", Score: 3},
			},
		},
		{
			Text: "我不够好，周围的人都比我优秀",
			Options: []Option{
				{Content: "确实", Score: 1},
				{Content: "有时", Score: 2},
				{Content: "不是", Score: 3},
			},
		},
		{
			Text: "我很清楚真正的自己是什么样的",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "我内心有真正追求的东西",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "我一定要不断往上爬、变得更厉害",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "外人的评价对我来说无所谓。",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "对象超过5小时没回消息，说自己窜稀了，你会怎么想？",
			Options: []Option{
				{Content: "拉稀不可能5小时，也许ta隐瞒了我。", Score: 1},
				{Content: "在信任和怀疑之间摇摆。", Score: 2},
				{Content: "也许今天ta真的不太舒服。", Score: 3},
			},
		},
		{
			Text: "我在感情里经常担心被对方抛弃",
			Options: []Option{
				{Content: "是的", Score: 1},
				{Content: "偶尔", Score: 2},
				{Content: "不是", Score: 3},
			},
		},
		{
			Text: "我对天发誓，我对待每一份感情都是认真的！",
			Options: []Option{
				{Content: "并没有", Score: 1},
				{Content: "也许？", Score: 2},
				{Content: "是的！（问心无愧骄傲脸）", Score: 3},
			},
		},
		{
			Text: "你的恋爱对象是一个尊老爱幼，温柔敦厚，洁身自好，光明磊落，大义凛然，能言善辩，口才流利，观察入微，见多识广，博学多才，诲人不倦，和蔼可亲，平易近人，心地善良，慈眉善目，积极进取，意气风发，玉树临风，国色天香，倾国倾城，花容月貌的人，此时你会？",
			Options: []Option{
				{Content: "就算ta再优秀我也不会陷入太深。", Score: 1},
				{Content: "会介于A和C之间。", Score: 2},
				{Content: "会非常珍惜ta，也许会变成恋爱脑。", Score: 3},
			},
		},
		{
			Text: "恋爱后，对象非常黏人，你作何感想？",
			Options: []Option{
				{Content: "那很爽了", Score: 1},
				{Content: "都行无所谓", Score: 2},
				{Content: "我更喜欢保留独立空间", Score: 3},
			},
		},
		{
			Text: "我在任何关系里都很重视个人空间",
			Options: []Option{
				{Content: "我更喜欢依赖与被依赖", Score: 1},
				{Content: "看情况", Score: 2},
				{Content: "是的！（斩钉截铁地说道）", Score: 3},
			},
		},
		{
			Text: "大多数人是善良的",
			Options: []Option{
				{Content: "其实邪恶的人心比世界上的痔疮更多。", Score: 1},
				{Content: "也许吧。", Score: 2},
				{Content: "是的，我愿相信好人更多。", Score: 3},
			},
		},
		{
			Text: "你走在街上，一位萌萌的小女孩蹦蹦跳跳地朝你走来（正脸、侧脸看都萌，用vivo、苹果、华为、OPPO手机看都萌，实在是非常萌的那种），她递给你一根棒棒糖，此时你作何感想？",
			Options: []Option{
				{Content: "呜呜她真好真可爱！居然给我棒棒糖！", Score: 3},
				{Content: "一脸懵逼，作挠头状", Score: 2},
				{Content: "这也许是一种新型诈骗？还是走开为好。", Score: 1},
			},
		},
		{
			Text: "快考试了，学校规定必须上晚自习，请假会扣分，但今晚你约了女/男神一起玩《绝地求生：刺激战场》（一款刺激的游戏），你怎么办？",
			Options: []Option{
				{Content: "翘了！反正就一次！", Score: 1},
				{Content: "干脆请个假吧。", Score: 2},
				{Content: "都快考试了还去啥。", Score: 3},
			},
		},
		{
			Text: "我喜欢打破常规，不喜欢被束缚",
			Options: []Option{
				{Content: "认同", Score: 1},
				{Content: "保持中立", Score: 2},
				{Content: "不认同", Score: 3},
			},
		},
		{
			Text: "我做事通常有目标。",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "突然某一天，我发现人生哪有什么真正的意义，某种程度，人不过是和动物一样被各种欲望支配着，被激素控制着，饿了就吃，困了就睡，只是一种无聊又低级的生物罢了，本质跟猫猫狗狗没什么区别。",
			Options: []Option{
				{Content: "是这样的。", Score: 1},
				{Content: "也许是，也许不是。", Score: 2},
				{Content: "这简直是胡扯", Score: 3},
			},
		},
		{
			Text: "我做事主要为了取得成果和进步，而不是避免麻烦和风险。",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "你因便秘坐在马桶上（已长达30分钟），拉不出很难受。此时你更像",
			Options: []Option{
				{Content: "再坐三十分钟看看，说不定就有了。", Score: 1},
				{Content: "用力拍打自己的屁股并说：“死屁股，快拉啊！”", Score: 2},
				{Content: "使用开塞露，快点拉出来才好。", Score: 3},
			},
		},
		{
			Text: "我做决定比较果断，不喜欢犹豫",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "此题没有题目，请盲选",
			Options: []Option{
				{Content: "反复思考后感觉应该选A？", Score: 1},
				{Content: "啊，要不选B？", Score: 2},
				{Content: "不会就选C？", Score: 3},
			},
		},
		{
			Text: "别人说你“执行力强”，你内心更接近哪句？",
			Options: []Option{
				{Content: "我被逼到最后确实执行力超强。。。", Score: 1},
				{Content: "啊，有时候吧。", Score: 2},
				{Content: "是的，事情本来就该被推进", Score: 3},
			},
		},
		{
			Text: "我做事常常有计划，____",
			Options: []Option{
				{Content: "然而计划不如变化快。", Score: 1},
				{Content: "有时能完成，有时不能。", Score: 2},
				{Content: "我讨厌被打破计划。", Score: 3},
			},
		},
		{
			Text: "你因玩《第五人格》（一款刺激的游戏）而结识许多网友，并被邀请线下见面，你的想法是？",
			Options: []Option{
				{Content: "网上口嗨下就算了，真见面还是有点忐忑。", Score: 1},
				{Content: "见网友也挺好，反正谁来聊我就聊两句。", Score: 2},
				{Content: "我会打扮一番并热情聊天，万一呢，我是说万一呢？", Score: 3},
			},
		},
		{
			Text: "朋友带了ta的朋友一起来玩，你最可能的状态是",
			Options: []Option{
				{Content: "对“朋友的朋友”天然有点距离感，怕影响二人关系", Score: 1},
				{Content: "看对方，能玩就玩。", Score: 2},
				{Content: "朋友的朋友应该也算我的朋友！要热情聊天", Score: 3},
			},
		},
		{
			Text: "我和人相处主打一个电子围栏，靠太近会自动报警。",
			Options: []Option{
				{Content: "认同", Score: 3},
				{Content: "中立", Score: 2},
				{Content: "不认同", Score: 1},
			},
		},
		{
			Text: "我渴望和我信任的人关系密切，熟得像失散多年的亲戚。",
			Options: []Option{
				{Content: "认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "不认同", Score: 3},
			},
		},
		{
			Text: "有时候你明明对一件事有不同的、负面的看法，但最后没说出来。多数情况下原因是：",
			Options: []Option{
				{Content: "这种情况较少。", Score: 1},
				{Content: "可能碍于情面或者关系。", Score: 2},
				{Content: "不想让别人知道自己是个阴暗的人。", Score: 3},
			},
		},
		{
			Text: "我在不同人面前会表现出不一样的自己",
			Options: []Option{
				{Content: "不认同", Score: 1},
				{Content: "中立", Score: 2},
				{Content: "认同", Score: 3},
			},
		},
		{
			Text: "您平时有什么爱好？",
			Options: []Option{
				{Content: "吃喝拉撒", Score: 1},
				{Content: "艺术爱好", Score: 2},
				{Content: "饮酒", Score: 3},
				{Content: "健身", Score: 4},
			},
		},
		{
			Text: "您对饮酒的态度是？",
			Options: []Option{
				{Content: "小酌怡情，喝不了太多。", Score: 1},
				{Content: "我习惯将白酒灌在保温杯，当白开水喝，酒精令我信服。", Score: 2},
			},
		},
	}

	SaveTestType(ctx, int(consts.SBTI), "SBTI傻大个性格测试", "SBTI傻大个性格测试基于十五维性格矩阵构建，使用曼哈顿距离算法将受测者的答题模式与预设人格模板进行精准匹配，融合现代心理测量学方法与趣味化设计理念，帮助用户以轻松方式发现自身独特性格特征与行为倾向，内含隐藏彩蛋设计，适合娱乐性自我探索场景。", []ReferenceItem{{Authors: []string{"Jung, C.G."}, Title: "Psychological Types", Journal: "Princeton University Press", Year: 1971}}, len(questions))

	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.SBTI).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.SBTI,
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

		optionStr := optionsToJSON(question.Options)

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
