package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

func PerceivedStress() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "当意想不到的事情发生时，我感到烦躁"},
		{Text: "我感到自己无法控制生活中重要的事情"},
		{Text: "我感到紧张和有压力"},
		{Text: "我对自己处理个人问题的能力缺乏信心"},
		{Text: "事情的发展不如我所期望"},
		{Text: "我感到无法处理所有必须做的事情"},
		{Text: "我难以控制生活中的恼怒情绪"},
		{Text: "我感到难以安排好自己的生活"},
		{Text: "因无法掌控发生的事情，我感到生气"},
		{Text: "困难积累得太大，我感到无法克服"},
		{Text: "我感到睡眠质量下降，难以入睡或易醒"},
		{Text: "我对平时喜欢的事情失去了兴趣"},
		{Text: "我感到身体疲惫，即使休息后也不易恢复"},
		{Text: "我难以找到有效的方法来放松自己"},
		{Text: "我感到头痛、胃痛或其他身体不适"},
		{Text: "我感到难以平衡工作与生活"},
		{Text: "我感到注意力难以集中，容易分心"},
		{Text: "我感到缺乏可以倾诉和寻求支持的人"},
		{Text: "我感到食欲改变（食欲不振或暴饮暴食）"},
		{Text: "我难以为自己设定合理的期望和目标"},
		{Text: "我感到心跳加快、呼吸急促或胸闷"},
		{Text: "我难以从日常活动中获得满足感和成就感"},
		{Text: "我容易对小事发脾气或感到不耐烦"},
		{Text: "我缺乏规律的体育锻炼习惯"},
		{Text: "我感到对未来缺乏希望和期待"},
	}

	options := []Option{
		{Tag: "A", Content: "从不", Score: 0},
		{Tag: "B", Content: "偶尔", Score: 1},
		{Tag: "C", Content: "有时", Score: 2},
		{Tag: "D", Content: "经常", Score: 3},
		{Tag: "E", Content: "总是", Score: 4},
	}
	// 转换为JSON
	optionStr := optionsToJSON(options)

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.PerceivedStress), "压力知觉", "感知压力量表(PSS)由Cohen、Kamarck和Mermelstein于1983年编制，通过14个条目评估个体在过去一个月内对生活事件不可预测性、不可控性及超负荷感的主观认知程度，不针对特定压力源，而是捕捉个体对整体生活压力的感知，是心理压力研究领域使用最广泛的标准化工具之一。", []ReferenceItem{{Authors: []string{"Cohen, S.", "Kamarck, T.", "Mermelstein, R."}, Title: "A global measure of perceived stress", Journal: "Journal of Health and Social Behavior", Year: 1983, DOI: "10.2307/2136404"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.PerceivedStress).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.PerceivedStress,
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
			log.Printf("插入选项失败: %v\n", err)
		}
	}
}
