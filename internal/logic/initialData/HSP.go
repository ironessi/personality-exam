package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// HSP 敏感度测试初始化
func HSP() {
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
		// 深度处理
		{Text: "我倾向于深入思考信息和经历"},
		{Text: "我对艺术和美学有强烈的反应"},
		{Text: "我能够理解他人的感受和情绪"},
		{Text: "我喜欢反思自己的行为和动机"},
		{Text: "我对音乐、文学和艺术有深刻的欣赏"},
		
		// 易过度刺激
		{Text: "我对噪音和强光敏感"},
		{Text: "我在社交场合后需要独处恢复能量"},
		{Text: "我容易被环境中的刺激所淹没"},
		{Text: "我对咖啡因和其他兴奋剂反应强烈"},
		{Text: "我在繁忙的环境中感到疲惫"},
		
		// 高情绪反应性
		{Text: "我对他人的情绪非常敏感"},
		{Text: "我容易被感动或流泪"},
		{Text: "我对他人的痛苦感到同情"},
		{Text: "我对批评和否定反应强烈"},
		{Text: "我能够感受到他人的情绪状态"},
		
		// 细微刺激觉察
		{Text: "我注意到环境中的细微变化"},
		{Text: "我能够察觉他人的微表情和非语言线索"},
		{Text: "我对气味和味道非常敏感"},
		{Text: "我注意到他人忽略的细节"},
		{Text: "我对温度和触觉刺激敏感"},
		
		// 其他特征
		{Text: "我在压力下容易感到不知所措"},
		{Text: "我喜欢安静和宁静的环境"},
		{Text: "我对药物和治疗的反应强烈"},
		{Text: "我倾向于避免暴力和恐怖的内容"},
		{Text: "我在做决定时考虑很多因素"},
		{Text: "我对生活中的美好事物有强烈的欣赏"},
		{Text: "我容易受到他人情绪的影响"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.HSP), "敏感度测试（HSP）", "高敏感人格量表(HSP)由Aron和Aron于1997年编制，基于感觉处理敏感性理论，评估个体对环境刺激、情绪信息和感官细节的深度处理倾向，约15-20%的人天生具有这种神经敏感特质，它并非缺陷而是一种独特的认知处理方式，与创造力、共情能力及深度思考密切相关。", []ReferenceItem{{Authors: []string{"Aron, E.N.", "Aron, A."}, Title: "Sensory-processing sensitivity and its relation to introversion and emotionality", Journal: "Journal of Personality and Social Psychology", Year: 1997, DOI: "10.1037/0022-3514.73.2.345"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.HSP,
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