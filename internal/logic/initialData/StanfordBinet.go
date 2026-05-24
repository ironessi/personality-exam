package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// StanfordBinet 斯坦福-比奈智力量表初始化
func StanfordBinet() {
	ctx := context.Background()

	var questions = []Question{
		{Text: "以下哪个数字是下一个序列中的数字：2, 4, 6, 8, ?", Options: []Option{{Tag: "A", Content: "10", Score: 4}, {Tag: "B", Content: "9", Score: 0}, {Tag: "C", Content: "12", Score: 0}, {Tag: "D", Content: "14", Score: 0}}},
		{Text: "'苹果'与'水果'的关系，就像'胡萝卜'与什么的关系？", Options: []Option{{Tag: "A", Content: "蔬菜", Score: 4}, {Tag: "B", Content: "植物", Score: 2}, {Tag: "C", Content: "食物", Score: 2}, {Tag: "D", Content: "橙色", Score: 0}}},
		{Text: "如果所有的猫都有尾巴，而咪咪是一只猫，那么咪咪有尾巴吗？", Options: []Option{{Tag: "A", Content: "是的", Score: 4}, {Tag: "B", Content: "不是", Score: 0}, {Tag: "C", Content: "可能", Score: 0}, {Tag: "D", Content: "不确定", Score: 0}}},
		{Text: "以下哪个词与其他三个不同类？", Options: []Option{{Tag: "A", Content: "红色", Score: 0}, {Tag: "B", Content: "蓝色", Score: 0}, {Tag: "C", Content: "绿色", Score: 0}, {Tag: "D", Content: "颜色", Score: 4}}},
		{Text: "5 + 3 × 2 = ?", Options: []Option{{Tag: "A", Content: "11", Score: 4}, {Tag: "B", Content: "16", Score: 0}, {Tag: "C", Content: "10", Score: 0}, {Tag: "D", Content: "13", Score: 0}}},
		{Text: "如果今天是星期一，那么后天是星期几？", Options: []Option{{Tag: "A", Content: "星期三", Score: 4}, {Tag: "B", Content: "星期二", Score: 0}, {Tag: "C", Content: "星期四", Score: 0}, {Tag: "D", Content: "星期五", Score: 0}}},
		{Text: "'勇敢'的反义词是什么？", Options: []Option{{Tag: "A", Content: "胆小", Score: 4}, {Tag: "B", Content: "坚强", Score: 0}, {Tag: "C", Content: "害怕", Score: 2}, {Tag: "D", Content: "懦弱", Score: 4}}},
		{Text: "将以下单词按逻辑顺序排列：1. 起床 2. 睡觉 3. 工作 4. 吃饭", Options: []Option{{Tag: "A", Content: "1-4-3-2", Score: 4}, {Tag: "B", Content: "1-3-4-2", Score: 2}, {Tag: "C", Content: "2-1-4-3", Score: 0}, {Tag: "D", Content: "3-4-1-2", Score: 0}}},
		{Text: "以下哪个图形与其他三个不同？", Options: []Option{{Tag: "A", Content: "圆形", Score: 0}, {Tag: "B", Content: "正方形", Score: 0}, {Tag: "C", Content: "三角形", Score: 0}, {Tag: "D", Content: "直线", Score: 4}}},
		{Text: "'书'是用来做什么的？", Options: []Option{{Tag: "A", Content: "阅读", Score: 4}, {Tag: "B", Content: "写字", Score: 0}, {Tag: "C", Content: "画画", Score: 0}, {Tag: "D", Content: "装饰", Score: 2}}},
	}

	// 保存测试类型信息
	SaveTestType(ctx, consts.StanfordBinet, "斯坦福-比奈智力量表", "斯坦福-比奈智力量表由Binet和Simon于1905年创制，经Terman于1916年修订为斯坦福版本，是世界上最早的标准化智力测验，通过评估语言理解、数量推理、抽象思维、视觉空间处理及工作记忆等多个认知维度，全面衡量个体的智力水平与认知能力结构。", []ReferenceItem{{Authors: []string{"Terman, L.M.", "Merrill, M.A."}, Title: "Stanford-Binet Intelligence Scale", Journal: "Houghton Mifflin", Year: 1960}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.StanfordBinet).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.StanfordBinet,
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
