package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"log"
)

// RIASEC 霍兰德职业兴趣测试初始化
func RIASEC() {
	ctx := context.Background()

	// 定义通用选项（5点计分，Tag A-E）
	options := []Option{
		{Tag: "A", Content: "非常不喜欢", Score: 1},
		{Tag: "B", Content: "不喜欢", Score: 2},
		{Tag: "C", Content: "一般", Score: 3},
		{Tag: "D", Content: "喜欢", Score: 4},
		{Tag: "E", Content: "非常喜欢", Score: 5},
	}
	optionStr := optionsToJSON(options)

	var questions = []Question{
		{Text: "我喜欢使用工具或机器来完成工作"},
		{Text: "我喜欢思考科学或数学问题"},
		{Text: "我喜欢创作艺术作品，如绘画、音乐或写作"},
		{Text: "我喜欢帮助他人解决问题"},
		{Text: "我喜欢领导或管理他人"},
		{Text: "我喜欢按照规则和程序工作"},
		{Text: "我喜欢修理电器或机械"},
		{Text: "我喜欢做科学实验"},
		{Text: "我喜欢设计或装饰物品"},
		{Text: "我喜欢参与志愿者活动"},
		{Text: "我喜欢说服他人接受我的观点"},
		{Text: "我喜欢整理文件或数据"},
		{Text: "我喜欢户外工作"},
		{Text: "我喜欢阅读科学书籍或文章"},
		{Text: "我喜欢表演或参加戏剧活动"},
		{Text: "我喜欢倾听他人的问题并提供建议"},
		{Text: "我喜欢销售产品或服务"},
		{Text: "我喜欢使用电脑进行数据处理"},
		{Text: "我喜欢从事体力劳动"},
		{Text: "我喜欢分析数据或解决复杂问题"},
		{Text: "我喜欢写作或写诗"},
		{Text: "我喜欢教育或培训他人"},
		{Text: "我喜欢组织活动或会议"},
		{Text: "我喜欢遵循详细的指示"},
		{Text: "我喜欢操作重型机械"},
		{Text: "我喜欢研究新事物"},
		{Text: "我喜欢弹奏乐器"},
		{Text: "我喜欢照顾他人的健康或福利"},
		{Text: "我喜欢谈判或调解争端"},
		{Text: "我喜欢记录或保存信息"},
		{Text: "我喜欢园艺或种植植物"},
		{Text: "我喜欢提出新的理论或假设"},
		{Text: "我喜欢绘画或素描"},
		{Text: "我喜欢在团队中工作"},
		{Text: "我喜欢冒险或尝试新事物"},
		{Text: "我喜欢检查细节或错误"},
		{Text: "我喜欢建造或修理东西"},
		{Text: "我喜欢数学或统计"},
		{Text: "我喜欢设计服装或时尚物品"},
		{Text: "我喜欢与人交往"},
		{Text: "我喜欢设定目标并努力实现"},
		{Text: "我喜欢使用表格或电子表格"},
		{Text: "我喜欢使用工具进行手工制作"},
		{Text: "我喜欢阅读哲学或心理学书籍"},
		{Text: "我喜欢写故事或小说"},
		{Text: "我喜欢帮助他人成长或发展"},
		{Text: "我喜欢竞争或赢得比赛"},
		{Text: "我喜欢按照计划或日程工作"},
		{Text: "我喜欢从事体育活动或锻炼"},
		{Text: "我喜欢观察或研究自然现象"},
		{Text: "我喜欢摄影或拍摄照片"},
		{Text: "我喜欢参与社区活动"},
		{Text: "我喜欢管理项目或团队"},
		{Text: "我喜欢计算或处理数字"},
		{Text: "我喜欢使用重型工具或设备"},
		{Text: "我喜欢解决逻辑谜题或问题"},
		{Text: "我喜欢创作音乐或歌曲"},
		{Text: "我喜欢倾听他人的感受"},
		{Text: "我喜欢冒险或创业"},
		{Text: "我喜欢整理或分类物品"},
		{Text: "我喜欢从事农业或林业工作"},
		{Text: "我喜欢进行研究或调查"},
		{Text: "我喜欢参与舞蹈或表演活动"},
		{Text: "我喜欢帮助他人解决冲突"},
		{Text: "我喜欢影响或说服他人"},
		{Text: "我喜欢填写表格或申请"},
		{Text: "我喜欢使用手工工具制作东西"},
		{Text: "我喜欢阅读科学期刊或论文"},
		{Text: "我喜欢设计网页或图形"},
		{Text: "我喜欢与不同类型的人交往"},
		{Text: "我喜欢推销或宣传产品"},
		{Text: "我喜欢使用数据库或信息系统"},
		{Text: "我喜欢修理家具或家居物品"},
		{Text: "我喜欢解决科学问题"},
		{Text: "我喜欢创作艺术品或手工艺品"},
		{Text: "我喜欢从事社会工作或咨询"},
		{Text: "我喜欢管理或监督他人"},
		{Text: "我喜欢核对或验证信息"},
		{Text: "我喜欢从事建筑或施工工作"},
		{Text: "我喜欢学习新的科学知识"},
		{Text: "我喜欢参与戏剧或电影制作"},
		{Text: "我喜欢帮助他人克服困难"},
		{Text: "我喜欢从事销售或营销工作"},
		{Text: "我喜欢整理或归档文件"},
	}

	// 插入或更新测试类型
	SaveTestType(ctx, int(consts.RIASEC), "霍兰德职业兴趣测试（RIASEC）", "霍兰德职业兴趣测试(RIASEC)由Holland于1959年提出，将职业兴趣与人格特征分为现实型(R)、研究型(I)、艺术型(A)、社会型(S)、企业型(E)、传统型(C)六种类型，通过匹配个人兴趣与职业环境，帮助个体找到最适合的职业方向，是全球职业规划与生涯咨询领域最具影响力的理论模型。", []ReferenceItem{{Authors: []string{"Holland, J.L."}, Title: "A theory of vocational choice", Journal: "Journal of Counseling Psychology", Year: 1959, DOI: "10.1037/h0040767"}}, len(questions))

	// 删除旧问题
	_, err := dao.Questions.Ctx(ctx).Where("test_type_id = ?", consts.RIASEC).Delete()
	if err != nil {
		log.Printf("删除数据失败: %v\n", err)
		return
	}

	// 插入问题和选项
	for idx, question := range questions {
		result, err := dao.Questions.Ctx(ctx).Data(&do.Questions{
			TestTypeId: consts.RIASEC,
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
