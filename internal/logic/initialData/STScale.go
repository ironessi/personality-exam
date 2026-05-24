package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// STScale 社会地位指数测试ST量表初始化
func STScale() {
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
		{Text: "我的收入水平在同龄人中处于较高水平"},
		{Text: "我的职业在社会上很受尊重"},
		{Text: "我拥有较高的学历"},
		{Text: "我拥有较多的社会资源和人脉关系"},
		{Text: "我居住在高档社区或地段"},
		{Text: "我的工作具有较高的社会影响力"},
		{Text: "我经常参与高端社交活动"},
		{Text: "我拥有较高的社会知名度"},
		{Text: "我的家庭背景较为显赫"},
		{Text: "我能够轻松获得优质的教育资源"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.STScale, "社会地位指数测试ST量表", "社会地位指数ST量表通过经济地位、职业地位、教育程度与社会资本四个维度的综合评估，客观量化个体在社会分层体系中的相对位置，帮助研究者和个体理解社会经济地位对心理健康、行为模式及生活机遇的深远影响，是社会学与心理学交叉研究的标准化测量工具。", []ReferenceItem{{Authors: []string{"Hollingshead, A.B."}, Title: "Four factor index of social status", Journal: "Yale University Department of Sociology", Year: 1975}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.STScale).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.STScale,
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
