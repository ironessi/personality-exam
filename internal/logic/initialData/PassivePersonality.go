package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// PassivePersonality 被动型人格测试初始化
func PassivePersonality() {
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
		{Text: "我在团队中通常会听从他人的决定，而不是主动提出自己的意见"},
		{Text: "我经常会为了避免冲突而妥协自己的立场"},
		{Text: "我很难拒绝他人的请求，即使这会给自己带来不便"},
		{Text: "我更倾向于等待机会，而不是主动争取"},
		{Text: "我在社交场合中通常会让他人主导对话"},
		{Text: "我经常会感到自己的需求被忽略，但又不敢表达"},
		{Text: "我更愿意跟随他人的计划，而不是自己制定计划"},
		{Text: "我在面对批评时通常会默默接受，而不是为自己辩护"},
		{Text: "我很难主动与陌生人建立联系"},
		{Text: "我经常会推迟做决定，希望他人能帮我做决定"},
		{Text: "我在工作中通常会完成分配的任务，但很少主动承担额外的责任"},
		{Text: "我更倾向于保持现状，而不是寻求改变"},
		{Text: "我很难表达自己的情绪，尤其是负面情绪"},
		{Text: "我经常会感到自己被他人支配，但又不敢反抗"},
		{Text: "我在选择时通常会询问他人的意见，而不是相信自己的判断"},
		{Text: "我很难设定个人边界，容易被他人侵犯"},
		{Text: "我更愿意在幕后工作，而不是成为焦点"},
		{Text: "我经常会感到自己的意见不重要"},
		{Text: "我很难主动结束一段关系，即使这段关系让我不开心"},
		{Text: "我在面对机会时通常会犹豫不决，担心自己无法胜任"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.PassivePersonality, "被动型人格测试", "被动型人格测试聚焦个体在决策制定、人际互动与责任承担等关键场景中的被动行为倾向，精准识别受测者的行为反应模式与被动倾向强度，帮助个体了解自身在人际关系、职场环境及日常生活中的互动特征，为个人成长与行为模式调整提供参考依据。", []ReferenceItem{{Authors: []string{"Millon, T."}, Title: "Disorders of Personality: DSM-IV and Beyond", Journal: "John Wiley & Sons", Year: 1996}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.PassivePersonality,
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
