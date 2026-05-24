package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// AMS 成就动机测试(AMS)初始化
func AMS() {
	ctx := context.Background()

	// 定义通用选项（添加Tags）
	options := []Option{
		{Tag: "A", Content: "非常不符合", Score: 1},
		{Tag: "B", Content: "不太符合", Score: 2},
		{Tag: "C", Content: "比较符合", Score: 3},
		{Tag: "D", Content: "非常符合", Score: 4},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我喜欢完成有挑战性的任务"},
		{Text: "我害怕在重要任务中失败"},
		{Text: "我喜欢设定高目标并努力实现"},
		{Text: "我担心自己会犯错误"},
		{Text: "我喜欢尝试新的方法来解决问题"},
		{Text: "我害怕被他人评价为不够优秀"},
		{Text: "我喜欢竞争并争取胜利"},
		{Text: "我担心失败会让我失去自信"},
		{Text: "我喜欢接受能够展示自己能力的任务"},
		{Text: "我害怕在他人面前表现不佳"},
		{Text: "我喜欢解决复杂的问题"},
		{Text: "我担心自己会让他人失望"},
		{Text: "我喜欢挑战自己的极限"},
		{Text: "我害怕失败会影响他人对我的看法"},
		{Text: "我喜欢完成需要付出努力的任务"},
		{Text: "我担心自己的能力不足"},
		{Text: "我喜欢设定具有挑战性的目标"},
		{Text: "我害怕在重要场合表现失常"},
		{Text: "我喜欢超越自己的过去表现"},
		{Text: "我担心自己会犯低级错误"},
		{Text: "我喜欢承担责任重大的任务"},
		{Text: "我害怕失败会让我感到羞耻"},
		{Text: "我喜欢在工作中获得认可和赞赏"},
		{Text: "我担心自己无法达到预期目标"},
		{Text: "我喜欢参与需要创造力的活动"},
		{Text: "我害怕在竞争中失利"},
		{Text: "我喜欢努力工作以实现自己的目标"},
		{Text: "我担心自己会被淘汰"},
		{Text: "我喜欢面对困难并克服它们"},
		{Text: "我害怕自己的努力会白费"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.AMS, "成就动机测试(AMS)", "成就动机量表(AMS)由挪威心理学家Gjesme和Nygard于1970年编制，通过追求成功的动机(MS)与回避失败的动机(MF)两个维度评估个体的成就动机结构，揭示个体在面对挑战性任务时的内在驱动力与心理防御机制，广泛应用于教育心理学、组织行为学与运动心理学研究。", []ReferenceItem{{Authors: []string{"Gjesme, T.", "Nygard, R."}, Title: "Achievement-related motives: Theoretical considerations and a test", Journal: "Scandinavian Journal of Educational Research", Year: 1970}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.AMS).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.AMS,
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
