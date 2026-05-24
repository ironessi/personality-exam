package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// SocialStatusIndex 社会地位指数测试ST量表初始化
func SocialStatusIndex() {
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
		{Text: "我对自己的收入水平感到满意"},
		{Text: "我拥有一定的资产规模"},
		{Text: "我有较强的消费能力"},
		{Text: "我的职业具有较高的社会认可度"},
		{Text: "我的职业技能要求较高"},
		{Text: "我的岗位权责范围较大"},
		{Text: "我拥有较高的教育程度"},
		{Text: "我的学历层次较高"},
		{Text: "我有广泛的社会关系网络"},
		{Text: "我的社会关系网络质量较高"},
		{Text: "我经常参与社会活动"},
		{Text: "我在社交场合中表现自信"},
		{Text: "我能轻松融入不同的社交圈子"},
		{Text: "我有较多的朋友和熟人"},
		{Text: "我的朋友和熟人在社会上有一定的影响力"},
		{Text: "我经常参加商务活动"},
		{Text: "我有较强的沟通能力"},
		{Text: "我善于建立和维护人际关系"},
		{Text: "我对社会热点问题有自己的见解"},
		{Text: "我关注社会发展趋势"},
		{Text: "我有较强的社会责任感"},
		{Text: "我经常参与公益活动"},
		{Text: "我对自己的社会地位感到满意"},
		{Text: "我认为自己在社会中扮演着重要的角色"},
		{Text: "我有较强的自我价值感"},
		{Text: "我对未来的社会地位提升充满信心"},
		{Text: "我善于抓住社会发展机遇"},
		{Text: "我有较强的适应社会变化的能力"},
		{Text: "我关注自己的社会形象"},
		{Text: "我注重自己的言行举止"},
		{Text: "我有较强的社会适应能力"},
		{Text: "我能与不同背景的人和谐相处"},
		{Text: "我有较强的团队合作能力"},
		{Text: "我善于协调人际关系"},
		{Text: "我有较强的领导能力"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.SocialStatusIndex, "社会地位指数测试ST量表", "社会地位指数量表通过经济地位、职业地位、教育程度与社会资本四个维度的综合评估，客观量化个体在社会分层体系中的相对位置，帮助研究者理解社会经济地位对心理健康、教育机遇与生活质量的影响，是社会学与心理学交叉研究中评估社会经济地位的标准化工具。", []ReferenceItem{{Authors: []string{"Hollingshead, A.B."}, Title: "Four factor index of social status", Journal: "Yale University Department of Sociology", Year: 1975}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.SocialStatusIndex,
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