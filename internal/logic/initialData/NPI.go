package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// NPI 自恋人格量表初始化
func NPI() {
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
		// 自我重要感
		{Text: "我认为自己是一个特别的人"},
		{Text: "我觉得自己比大多数人更有价值"},
		{Text: "我相信自己有特殊的天赋和能力"},
		{Text: "我认为自己应该得到特别的待遇"},
		
		// 对成功和权力的幻想
		{Text: "我经常幻想自己的成功和成就"},
		{Text: "我渴望获得权力和影响力"},
		{Text: "我相信自己有一天会成为重要人物"},
		{Text: "我对名望和地位有强烈的渴望"},
		
		// 独特性和特殊性
		{Text: "我认为自己的想法和观点很特别"},
		{Text: "我觉得自己与普通人不同"},
		{Text: "我需要他人认识到我的独特性"},
		{Text: "我相信自己有特殊的使命或命运"},
		
		// 需要过度赞美
		{Text: "我喜欢被他人赞美和钦佩"},
		{Text: "我需要他人的认可和关注"},
		{Text: "我喜欢成为众人关注的焦点"},
		{Text: "我对他人的评价很在意"},
		
		// 权利感
		{Text: "我认为自己应该得到最好的"},
		{Text: "我觉得规则对我来说不那么重要"},
		{Text: "我期望他人按照我的意愿行事"},
		{Text: "我对他人的要求很高"},
		
		// 人际剥削
		{Text: "我会利用他人来达到自己的目的"},
		{Text: "我觉得自己有权利使用他人"},
		{Text: "我不介意操纵他人以获得好处"},
		{Text: "我认为为了成功可以不择手段"},
		
		// 缺乏同理心
		{Text: "我不太关心他人的感受"},
		{Text: "我觉得自己的需求比他人更重要"},
		{Text: "我很难理解他人的观点"},
		{Text: "我不太容易同情他人的困境"},
		
		// 嫉妒和被嫉妒
		{Text: "我经常嫉妒他人的成功"},
		{Text: "我认为他人嫉妒我的成就"},
		{Text: "我对他人的好运感到不满"},
		{Text: "我觉得他人应该羡慕我的生活"},
		
		// 傲慢和自大
		{Text: "我认为自己比他人更聪明"},
		{Text: "我对他人的能力持怀疑态度"},
		{Text: "我觉得自己的判断比他人更准确"},
		{Text: "我对他人的缺点很敏感"},
		
		// 其他特征
		{Text: "我喜欢控制局面和他人"},
		{Text: "我对批评和否定反应强烈"},
		{Text: "我喜欢展示自己的成就和 possessions"},
		{Text: "我认为自己值得被特别对待"},
	}

	// 保存测试类型信息
	SaveTestType(ctx, int(consts.NPI), "自恋人格量表（NPI）", "自恋人格量表(NPI)由Raskin和Hall于1979年编制，包含40个强迫选择题，评估非临床人群中的自恋特质，测量自我重要感、权力幻想、独特性、需要赞美、权利感、人际剥削及缺乏同理心等维度，是人格心理学中自恋研究领域最广泛使用的标准化测量工具。", []ReferenceItem{{Authors: []string{"Raskin, R.", "Hall, C.S."}, Title: "A narcissistic personality inventory", Journal: "Psychological Reports", Year: 1979, DOI: "10.2466/pr0.1979.45.2.590"}}, len(questions))

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.NPI,
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