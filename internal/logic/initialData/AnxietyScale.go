package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// AnxietyScale 焦虑性测试 As量表初始化
func AnxietyScale() {
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
		{Text: "我感到紧张或焦虑"},
		{Text: "我感到心跳加速"},
		{Text: "我感到呼吸急促"},
		{Text: "我感到手心出汗"},
		{Text: "我感到头晕或眩晕"},
		{Text: "我感到胃部不适"},
		{Text: "我感到肌肉紧张"},
		{Text: "我感到颤抖或发抖"},
		{Text: "我感到口干舌燥"},
		{Text: "我感到喉咙发紧"},
		{Text: "我感到呼吸困难"},
		{Text: "我感到胸痛或胸闷"},
		{Text: "我感到恶心"},
		{Text: "我感到想上厕所"},
		{Text: "我感到疲惫或虚弱"},
		{Text: "我感到睡眠困难"},
		{Text: "我感到难以集中注意力"},
		{Text: "我感到记忆力下降"},
		{Text: "我感到易怒或烦躁"},
		{Text: "我感到担忧或不安"},
		{Text: "我感到害怕或恐惧"},
		{Text: "我感到惊慌失措"},
		{Text: "我感到失控"},
		{Text: "我感到即将崩溃"},
		{Text: "我感到死亡临近"},
		{Text: "我感到不真实或脱离现实"},
		{Text: "我感到周围环境不真实"},
		{Text: "我感到自我意识增强"},
		{Text: "我感到过度警觉"},
		{Text: "我感到容易受惊"},
		{Text: "我感到避免某些情境"},
		{Text: "我感到需要逃离某些情境"},
		{Text: "我感到需要寻求安全感"},
		{Text: "我感到需要有人陪伴"},
		{Text: "我感到无法放松"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.AnxietyScale), "焦虑性测试 As量表", "焦虑量表(As量表)由Welsh于1956年从明尼苏达多项人格测验(MMPI)中提取编制，专门评估个体主观感受到的焦虑情绪严重程度及其伴随的生理紧张、心理不安等症状，能够有效区分特质焦虑与状态焦虑，是临床焦虑评估与心理健康筛查的经典标准化工具。", []ReferenceItem{{Authors: []string{"Welsh, G.S."}, Title: "Factor dimensions A and R", Journal: "Basic readings on the MMPI in psychology and medicine", Year: 1956}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.AnxietyScale,
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
