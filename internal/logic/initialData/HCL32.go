package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// HCL32 轻躁狂症状自评量表初始化
func HCL32() {
	ctx := context.Background()

	// 定义选项（HCL-32使用是/否计分）
	options := []Option{
		{Tag: "A", Content: "是", Score: 1},
		{Tag: "B", Content: "否", Score: 0},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "你是否感到异常快乐、兴奋或乐观？"},
		{Text: "你是否比平时更加活跃、精力充沛？"},
		{Text: "你是否感觉睡很少的时间也能精神饱满？"},
		{Text: "你是否感觉脑子里有很多想法，思维跳跃很快？"},
		{Text: "你是否比平时更喜欢与人交往？"},
		{Text: "你是否比平时更爱说话？"},
		{Text: "你是否比平时更自信？"},
		{Text: "你是否比平时更冲动？"},
		{Text: "你是否比平时更容易感到烦躁或易怒？"},
		{Text: "你是否比平时更有创造力？"},
		{Text: "你是否比平时更爱冒险？"},
		{Text: "你是否比平时更关注自己的外貌或穿着？"},
		{Text: "你是否比平时更健谈？"},
		{Text: "你是否比平时更容易分心？"},
		{Text: "你是否比平时更有性欲？"},
		{Text: "你是否比平时更爱花钱？"},
		{Text: "你是否比平时更爱社交？"},
		{Text: "你是否比平时更爱旅行？"},
		{Text: "你是否比平时更爱尝试新事物？"},
		{Text: "你是否比平时更爱做计划？"},
		{Text: "你是否比平时更爱表现自己？"},
		{Text: "你是否比平时更爱争论？"},
		{Text: "你是否比平时更爱挑战权威？"},
		{Text: "你是否比平时更爱冒险？"},
		{Text: "你是否比平时更爱赌博？"},
		{Text: "你是否比平时更爱喝酒？"},
		{Text: "你是否比平时更爱吸烟？"},
		{Text: "你是否比平时更爱吃东西？"},
		{Text: "你是否比平时更爱运动？"},
		{Text: "你是否比平时更爱工作？"},
		{Text: "你是否比平时更爱学习？"},
		{Text: "你是否比平时更爱思考？"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.HCL32), "轻躁狂症状自评量表(HCL-32)", "轻躁狂症状自评量表(HCL-32)由Angst等于2005年编制，包含32个条目，专门用于识别双相情感障碍（尤其是II型双相障碍）患者的轻躁狂发作史，通过评估精力增加、睡眠减少、思维加速等轻躁狂特征症状，帮助临床医生区分单相抑郁与双相抑郁，是双相障碍早期筛查的重要工具。", []ReferenceItem{{Authors: []string{"Angst, J.", "Adolfsson, R.", "Benazzi, F."}, Title: "The HCL-32: Towards a self-assessment tool for hypomanic symptoms in outpatients", Journal: "Journal of Affective Disorders", Year: 2005, DOI: "10.1016/j.jad.2005.05.011"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.HCL32,
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