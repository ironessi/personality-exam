package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func SocialDistress() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "面对一群人时，我会觉得自己不知道该把手放在哪里，特别局促。"},
		{Text: "在社交场合结束后，我常常会回想自己有没有哪里做得不好，然后很懊恼。"},
		{Text: "我很少主动在社交平台上和不太熟悉的人聊天互动。"},
		{Text: "如果要去一个新的社交圈子，我会特别害怕融入不进去。"},
		{Text: "我觉得自己在社交中总是很笨拙，不像别人那样自然。"},
		{Text: "和别人聊天时，我常常担心冷场，所以感觉压力很大。"},
		{Text: "我不太愿意在别人面前展示自己的才艺或特长，哪怕别人鼓励我也不想。"},
		{Text: "在社交场合中，别人的一个眼神或者表情变化，都会让我想很多，心里不舒服。"},
		{Text: "我会尽量避免参加那种需要和很多人互动交流的会议。"},
		{Text: "我害怕成为众人关注的焦点，哪怕是因为好的事情。"},
		{Text: "我觉得和别人相处时，很难找到合适的话题，这让我很烦恼。"},
		{Text: "在社交场合中，我很难放松下来，总是紧绷着神经。"},
		{Text: "当别人邀请我参加社交活动时，我第一反应往往是想拒绝。"},
		{Text: "我不太喜欢参与有很多人参与的游戏或活动。"},
		{Text: "即使在熟悉的人群中，我也常常不愿主动开口说话。"},
		{Text: "我对主动邀请别人一起做事情感到很为难。"},
		{Text: "在集体讨论中，我总是担心自己的发言不够好，所以不太敢说话。"},
		{Text: "如果可以选择，我宁愿一个人待着也不愿意参加社交活动。"},
		{Text: "我害怕在社交场合中说错话被别人笑话。"},
		{Text: "和不太熟悉的人聊天时，我会觉得很不自在。"},
		{Text: "在社交场合，我会很在意别人对我的看法，这让我很苦恼。"},
		{Text: "我常常找借口推脱参加朋友组织的聚会活动。"},
		{Text: "想到要在众人面前做个简单的自我介绍，我就会感到不安。"},
		{Text: "在社交场合中，如果别人不主动和我说话，我基本不会先开口。"},
		{Text: "我很少主动去结识新的朋友。"},
		{Text: "当不得不进入一个满是陌生人的场合时，我心里会很紧张。"},
		{Text: "我尽量避免和陌生人有眼神接触。"},
		{Text: "参加社交聚会时，我总是会提前很久就开始担心。"},
	}

	// 定义选项
	options := []Option{
		{Tag: "A", Content: "是", Score: 0},
		{Tag: "B", Content: "否", Score: 1},
	}

	// 转换为JSON
	optionStr := optionsToJSON(options)

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.SocialDistress), "社会痛苦", "社会回避及苦恼量表(SAD)由Watson和Friend于1969年编制，分别测量个体在社交场合中的主动回避行为倾向与主观痛苦体验两个独立维度，能够区分因焦虑而回避社交与因不感兴趣而回避社交的不同模式，是评估社交焦虑与社交回避行为的经典标准化工具。", []ReferenceItem{{Authors: []string{"Watson, D.", "Friend, R."}, Title: "Measurement of social-evaluative anxiety", Journal: "Journal of Consulting and Clinical Psychology", Year: 1969, DOI: "10.1037/h0027806"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.SocialDistress).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.SocialDistress,
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
			log.Printf("插入选项失败 [%d]: %v\n", idx+1, err)
			break
		}
	}
}
