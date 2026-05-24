package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// Rorschach 罗夏墨迹测验初始化
func Rorschach() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "你从这幅墨迹图片中看到了什么？（描述你看到的内容、动作、情绪等）", Options: []Option{{Tag: "A", Content: "动物或动物的一部分", Score: 1}, {Tag: "B", Content: "人物或人物的一部分", Score: 2}, {Tag: "C", Content: "物体或物体的一部分", Score: 3}, {Tag: "D", Content: "抽象的形状或图案", Score: 4}}},
		{Text: "你看到的东西在做什么？", Options: []Option{{Tag: "A", Content: "在运动或互动", Score: 1}, {Tag: "B", Content: "静止不动", Score: 2}, {Tag: "C", Content: "在变化或变形", Score: 3}, {Tag: "D", Content: "没有特定的动作", Score: 4}}},
		{Text: "这幅图片让你联想到了什么情绪？", Options: []Option{{Tag: "A", Content: "积极的情绪（快乐、喜悦等）", Score: 1}, {Tag: "B", Content: "中性的情绪（平静、无情绪等）", Score: 2}, {Tag: "C", Content: "消极的情绪（恐惧、愤怒等）", Score: 3}, {Tag: "D", Content: "复杂的情绪组合", Score: 4}}},
		{Text: "你是从图片的哪个部分看到这个内容的？", Options: []Option{{Tag: "A", Content: "整个图片", Score: 1}, {Tag: "B", Content: "主要部分", Score: 2}, {Tag: "C", Content: "次要部分", Score: 3}, {Tag: "D", Content: "细节部分", Score: 4}}},
		{Text: "是什么让你觉得这像你所描述的内容？", Options: []Option{{Tag: "A", Content: "形状", Score: 1}, {Tag: "B", Content: "颜色", Score: 2}, {Tag: "C", Content: "纹理", Score: 3}, {Tag: "D", Content: "整体印象", Score: 4}}},
		{Text: "这幅图片让你想起了什么具体的事物或情境？", Options: []Option{{Tag: "A", Content: "具体的物体", Score: 1}, {Tag: "B", Content: "具体的场景", Score: 2}, {Tag: "C", Content: "抽象的概念", Score: 3}, {Tag: "D", Content: "没有特定的联想", Score: 4}}},
		{Text: "你看到的内容是清晰明确的还是模糊不确定的？", Options: []Option{{Tag: "A", Content: "非常清晰", Score: 1}, {Tag: "B", Content: "比较清晰", Score: 2}, {Tag: "C", Content: "比较模糊", Score: 3}, {Tag: "D", Content: "非常模糊", Score: 4}}},
		{Text: "这幅图片给你的整体感觉是什么？", Options: []Option{{Tag: "A", Content: "愉快的", Score: 1}, {Tag: "B", Content: "中性的", Score: 2}, {Tag: "C", Content: "不愉快的", Score: 3}, {Tag: "D", Content: "复杂的", Score: 4}}},
		{Text: "你看到的内容是有生命的还是无生命的？", Options: []Option{{Tag: "A", Content: "有生命的", Score: 1}, {Tag: "B", Content: "无生命的", Score: 2}, {Tag: "C", Content: "两者都有", Score: 3}, {Tag: "D", Content: "不确定", Score: 4}}},
		{Text: "这幅图片让你产生了什么身体感觉？", Options: []Option{{Tag: "A", Content: "没有特别的感觉", Score: 1}, {Tag: "B", Content: "舒适的感觉", Score: 2}, {Tag: "C", Content: "不舒服的感觉", Score: 3}, {Tag: "D", Content: "强烈的感觉", Score: 4}}},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.Rorschach, "罗夏墨迹测验（Rorschach）", "罗夏墨迹测验由瑞士精神科医生Hermann Rorschach于1921年创制，通过呈现10张标准化墨迹图片，分析个体的自由联想反应来评估人格结构、情绪状态、思维过程与潜意识动机，是投射测验的经典代表，广泛应用于临床心理诊断、人格评估与心理治疗领域。", []ReferenceItem{{Authors: []string{"Rorschach, H."}, Title: "Psychodiagnostik", Journal: "Ernst Bircher", Year: 1921}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.Rorschach).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.Rorschach,
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
