package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// Zootopia 疯狂动物城测试性格初始化
func Zootopia() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "在团队中，我更倾向于：", Options: []Option{{Tag: "A", Content: "提出新想法和创意", Score: 1}, {Tag: "B", Content: "组织和协调团队成员", Score: 2}, {Tag: "C", Content: "执行具体任务", Score: 3}, {Tag: "D", Content: "提供支持和鼓励", Score: 4}}},
		{Text: "面对问题，我通常会：", Options: []Option{{Tag: "A", Content: "直接面对，迅速解决", Score: 1}, {Tag: "B", Content: "仔细分析，制定计划", Score: 2}, {Tag: "C", Content: "寻求他人的帮助", Score: 3}, {Tag: "D", Content: "等待最佳时机", Score: 4}}},
		{Text: "我更喜欢的生活方式是：", Options: []Option{{Tag: "A", Content: "充满变化和惊喜", Score: 1}, {Tag: "B", Content: "有规律和计划", Score: 2}, {Tag: "C", Content: "舒适和放松", Score: 3}, {Tag: "D", Content: "忙碌和充实", Score: 4}}},
		{Text: "当我遇到新环境时，我会：", Options: []Option{{Tag: "A", Content: "立即探索和适应", Score: 1}, {Tag: "B", Content: "观察后再行动", Score: 2}, {Tag: "C", Content: "寻求熟悉的人帮助", Score: 3}, {Tag: "D", Content: "感到紧张和不安", Score: 4}}},
		{Text: "我做决定的方式是：", Options: []Option{{Tag: "A", Content: "凭直觉和感觉", Score: 1}, {Tag: "B", Content: "理性分析和逻辑", Score: 2}, {Tag: "C", Content: "参考他人意见", Score: 3}, {Tag: "D", Content: "犹豫不决", Score: 4}}},
		{Text: "我在团队中的角色通常是：", Options: []Option{{Tag: "A", Content: "领导者", Score: 1}, {Tag: "B", Content: "策划者", Score: 2}, {Tag: "C", Content: "执行者", Score: 3}, {Tag: "D", Content: "支持者", Score: 4}}},
		{Text: "面对挫折，我会：", Options: []Option{{Tag: "A", Content: "积极寻找解决方案", Score: 1}, {Tag: "B", Content: "冷静分析原因", Score: 2}, {Tag: "C", Content: "寻求他人支持", Score: 3}, {Tag: "D", Content: "暂时退缩", Score: 4}}},
		{Text: "我更喜欢的娱乐方式是：", Options: []Option{{Tag: "A", Content: "参加派对或社交活动", Score: 1}, {Tag: "B", Content: "阅读或观看电影", Score: 2}, {Tag: "C", Content: "户外运动", Score: 3}, {Tag: "D", Content: "宅在家里休息", Score: 4}}},
		{Text: "我对规则的态度是：", Options: []Option{{Tag: "A", Content: "灵活变通", Score: 1}, {Tag: "B", Content: "严格遵守", Score: 2}, {Tag: "C", Content: "视情况而定", Score: 3}, {Tag: "D", Content: "不太在意", Score: 4}}},
		{Text: "我与他人沟通的方式是：", Options: []Option{{Tag: "A", Content: "直接坦率", Score: 1}, {Tag: "B", Content: "委婉含蓄", Score: 2}, {Tag: "C", Content: "倾听多于表达", Score: 3}, {Tag: "D", Content: "根据对方调整", Score: 4}}},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.Zootopia, "疯狂动物城测试性格", "疯狂动物城性格测试以电影中16位核心角色为原型，融合MBTI人格理论与社会心理学视角，通过40道情境选择题将受测者与对应的动漫角色进行性格匹配，以轻松有趣的方式帮助用户了解自身性格特征与行为模式，适合自我探索与朋友间的互动分享。", []ReferenceItem{{Authors: []string{"Myers, I.B.", "McCaulley, M.H."}, Title: "Manual: A guide to the development and use of the Myers-Briggs Type Indicator", Journal: "Consulting Psychologists Press", Year: 1985}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.Zootopia).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.Zootopia,
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

		// 将选项转换为JSON
		optionStr := optionsToJSON(question.Options)

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
