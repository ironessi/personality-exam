package initialData

import (
	"PersonalityTest/internal/consts"
	"PersonalityTest/internal/dao"
	"context"
	"encoding/json"
	"runtime/debug"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type Option struct {
	Tag     string `json:"tag" dc:"选项标签"`
	Content string `json:"content" dc:"选项内容"`
	Score   int    `json:"score" dc:"选项分数"`
}

type ReferenceItem struct {
	Authors []string `json:"authors"`
	Title   string   `json:"title"`
	Journal string   `json:"journal"`
	Year    int      `json:"year"`
	DOI     string   `json:"doi,omitempty"`
}

type Question struct {
	Text       string
	Options    []Option
	OptionText string
}

// ImportData 智能初始化数据 - 总是强制重新初始化
func ImportData() {
	ctx := context.Background()

	g.Log().Info(ctx, "开始初始化数据...")
	g.Log().Info(ctx, "强制重新初始化所有测试数据，不管现有数据情况")

	// 总是调用ForceImportData，确保实际的问卷题目被插入到数据库中
	ForceImportData()

	// 验证最终结果
	finalQuestionCount, _ := dao.Questions.Ctx(ctx).Count()
	g.Log().Info(ctx, "数据初始化完成，最终问题数量: ", finalQuestionCount)
}

// ForceImportData 强制重新初始化所有数据，但保留用户测试记录
func ForceImportData() {
	ctx := context.Background()

	g.Log().Warning(ctx, "⚠️ 强制重新初始化测试数据...")

	// 在初始化之前，先查询用户测试记录的数量，用于验证保留是否成功
	beforeUserTestCount, _ := dao.UserTestResults.Ctx(ctx).Count()
	beforeUserAnswerCount, _ := dao.UserAnswers.Ctx(ctx).Count()
	g.Log().Info(ctx, "初始化前用户测试记录数量:", beforeUserTestCount)
	g.Log().Info(ctx, "初始化前用户答案数量:", beforeUserAnswerCount)

	// 不再删除任何表的数据，避免级联删除用户测试记录
	// 直接插入或更新测试类型、题目和选项数据
	g.Log().Info(ctx, "✅ 保留所有现有数据，直接插入或更新测试类型、题目和选项")

	// 重新初始化 - 使用更健壮的方式
	g.Log().Info(ctx, "开始重新初始化...")
	initAllTestsWithRetry()

	// 验证结果
	finalCount, _ := dao.TestTypes.Ctx(ctx).Count()
	if finalCount >= 11 {
		g.Log().Info(ctx, "✅ 强制重新初始化完成，共添加", finalCount, "种测试类型")
	} else {
		g.Log().Error(ctx, "❌ 初始化不完整，只有", finalCount, "种测试类型")
	}

	// 在初始化之后，再查询用户测试记录的数量，用于验证保留是否成功
	afterUserTestCount, _ := dao.UserTestResults.Ctx(ctx).Count()
	afterUserAnswerCount, _ := dao.UserAnswers.Ctx(ctx).Count()
	g.Log().Info(ctx, "初始化后用户测试记录数量:", afterUserTestCount)
	g.Log().Info(ctx, "初始化后用户答案数量:", afterUserAnswerCount)

	// 比较初始化前后的数量，验证是否成功保留
	if afterUserTestCount == beforeUserTestCount {
		g.Log().Info(ctx, "✅ 成功保留用户测试记录数据")
	} else {
		g.Log().Warning(ctx, "⚠️ 用户测试记录数据可能丢失，初始化前:", beforeUserTestCount, "初始化后:", afterUserTestCount)
	}

	if afterUserAnswerCount == beforeUserAnswerCount {
		g.Log().Info(ctx, "✅ 成功保留用户答案数据")
	} else {
		g.Log().Warning(ctx, "⚠️ 用户答案数据可能丢失，初始化前:", beforeUserAnswerCount, "初始化后:", afterUserAnswerCount)
	}
}

// initAllTests 初始化所有测试（私有函数，被 ImportData 和 ForceImportData 调用）
// initAllTestsWithRetry 带重试的初始化
func initAllTestsWithRetry() {
	ctx := context.Background()

	initFunctions := []struct {
		id   int
		name string
		fn   func()
	}{
		{1, "SCL90", SCL90},
		{2, "生活满意度", GeneralLifeSatisfaction},
		{3, "感知压力", PerceivedStress},
		{4, "PROMIS焦虑", PROMISAnxiety},
		{5, "社交痛苦", SocialDistress},
		{6, "UCLA孤独", UCLALoneliness},
		{7, "生命意义感", LifeMeaning},
		{8, "PANAS情绪", Panas},
		{9, "自我效能感", SelfEfficacy},
		{10, "身体攻击性", PhysicalAggression},
		{11, "敌意倾向", Hostility},
		{12, "霍兰德职业兴趣测试", RIASEC},
		{13, "16型人格测试", MBTI},
		{14, "Personalities人格测试", Personalities},
		{15, "PSTRI压力知觉测试量表", PSTRI},
		{16, "菲尔人格十项测试题", PhilPersonality},
		{17, "社会地位指数测试ST量表", STScale},
		{18, "成就动机测试(AMS)", AMS},
		{19, "疯狂动物城测试性格", Zootopia},
		{20, "斯坦福-比奈智力量表", StanfordBinet},
		{21, "罗夏墨迹测验", Rorschach},
		// 新增心理测评问卷
		{22, "贝克抑郁自评量表", BDI},
		{23, "焦虑症自测量表(SDS)", SDSAnxiety},
		{24, "轻躁狂症状自评量表(HCL-32)", HCL32},
		{25, "心理健康素养量表", MHLAQ},
		{26, "29度人格测试", TwentyNineDegreePersonality},
		{27, "同理心通情量表", EMP},
		{28, "DISC性格测试", DISC},
		{29, "EPQ艾森克人格测试-国际标准版", EPQ},
		{30, "社会地位指数测试ST量表", SocialStatusIndex},
		{31, "马基雅维利人格测试", Machiavellianism},
		// 新增心理测评问卷
		{32, "自我控力指数评估", SelfControlIndex},
		{33, "被动型人格测试", PassivePersonality},
		{34, "色彩性格测试（True Colors）", TrueColors},
		{35, "服从性测试", ObedienceTest},
		{36, "多相人格测试", MMPI},
		{37, "焦虑性测试 As量表", AnxietyScale},
		{38, "爱德华个人偏好测验", EPPS},
		{39, "敏感度测试（HSP）", HSP},
		{40, "自恋人格量表（NPI）", NPI},
		{41, "贝姆性别角色测试量表", BemSexRoleInventory},
		{42, "SBTI傻大个性格测试", SBTI},
	}

	g.Log().Info(ctx, "开始初始化测试数据，共", len(initFunctions), "种测试")

	// 记录失败的测试
	var failedTests []string

	for _, item := range initFunctions {
		g.Log().Infof(ctx, "初始化第 %d/%d 种测试: %s...", item.id, len(initFunctions), item.name)

		startTime := time.Now()

		// 尝试3次
		success := false
		for attempt := 1; attempt <= 3; attempt++ {
			if attempt > 1 {
				g.Log().Infof(ctx, "  第 %d 次重试...", attempt)
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						g.Log().Errorf(ctx, "初始化 %s 时发生panic (尝试%d): %v", item.name, attempt, r)
						g.Log().Errorf(ctx, "panic堆栈: %s", debug.Stack())
					}
				}()

				// 初始化前先清理该测试的所有题目和选项，但保留测试类型
				// 这样可以避免重复键错误，同时确保数据的一致性
				g.Log().Infof(ctx, "  清理测试 %s 的现有题目和选项...", item.name)

				// 初始化前先清理该测试的所有题目和选项，但保留测试类型
				// 这样可以避免重复键错误，同时确保数据的一致性
				g.Log().Infof(ctx, "  清理测试 %s 的现有题目和选项...", item.name)

				// 直接使用子查询删除选项，避免Scan映射问题
				_, err := dao.Options.Ctx(ctx).
					Where("question_id IN (SELECT id FROM questions WHERE test_type_id = ?)", item.id).
					Delete()
				if err != nil {
					g.Log().Errorf(ctx, "  删除测试 %s 的选项失败: %v", item.name, err)
					return
				}

				// 删除该测试的所有题目
				_, err = dao.Questions.Ctx(ctx).
					Where("test_type_id", item.id).
					Delete()
				if err != nil {
					g.Log().Errorf(ctx, "  删除测试 %s 的题目失败: %v", item.name, err)
					return
				}

				g.Log().Infof(ctx, "  开始执行测试 %s 的初始化函数...", item.name)
				// 执行测试初始化函数
				item.fn()
				success = true
			}()

			if success {
				break
			}

			// 等待一下再重试
			time.Sleep(100 * time.Millisecond)
		}

		elapsed := time.Since(startTime)

		if success {
			g.Log().Infof(ctx, "✅ 第 %d 种测试 %s 初始化完成，耗时 %v", item.id, item.name, elapsed)

			// 立即验证这个测试类型的数据
			count, _ := dao.TestTypes.Ctx(ctx).Where("id", item.id).Count()
			questionCount, _ := dao.Questions.Ctx(ctx).Where("test_type_id", item.id).Count()
			g.Log().Infof(ctx, "   状态: 类型存在=%v, 问题数=%d", count > 0, questionCount)

			if count == 0 || questionCount == 0 {
				g.Log().Warningf(ctx, "   ⚠️  %s 数据可能不完整", item.name)
			}
		} else {
			g.Log().Errorf(ctx, "❌ 第 %d 种测试 %s 初始化失败", item.id, item.name)
			failedTests = append(failedTests, item.name)
		}
	}

	g.Log().Info(ctx, "所有测试初始化完成")

	if len(failedTests) > 0 {
		g.Log().Errorf(ctx, "以下测试初始化失败: %v", failedTests)
	}

	// 继续执行分类更新

	// 统一更新测试分类
	g.Log().Info(ctx, "开始更新测试分类信息...")
	categoryMap := map[int][]int{
		1: {1, 3, 4, 15, 22, 23, 24, 25, 37},                                // 健康测试
		2: {10, 11, 13, 14, 16, 21, 26, 28, 29, 31, 33, 34, 36, 39, 40, 41}, // 性格测试
		3: {2, 5, 6, 7, 8, 27},                                              // 情感测试
		4: {12},                                                             // 职业测试
		5: {19, 42},                                                         // 趣味测试
		6: {20},                                                             // 智商测试
		7: {9, 17, 18, 30, 32, 35, 38},                                      // 能力测试
	}

	for catID, testIDs := range categoryMap {
		_, err := dao.TestTypes.Ctx(ctx).
			Data(g.Map{"category_id": catID}).
			Where("id IN (?)", testIDs).
			Update()
		if err != nil {
			g.Log().Errorf(ctx, "更新分类 %d 失败: %v", catID, err)
		}
	}
	g.Log().Info(ctx, "测试分类信息更新完成")
}

// SaveTestType 统一保存测试类型信息，包含价格逻辑和基础元数据
func SaveTestType(ctx context.Context, id int, name, desc string, refs []ReferenceItem, questionCount int) {
	// 1. 基础定价规则
	price := 3.99
	if questionCount >= 25 && questionCount <= 50 {
		price = 4.99
	} else if questionCount > 50 {
		price = 6.99
	}

	// 2. 特殊量表（金牌量表）定价 9.99
	premiumTests := map[int]bool{
		consts.MBTI:          true,
		consts.MMPI:          true,
		consts.DISC:          true,
		consts.RIASEC:        true,
		consts.Personalities: true,
	}
	if premiumTests[id] {
		price = 9.99
	}

	// 3. 序列化 references 为 JSON 字符串
	refsJSON, err := json.Marshal(refs)
	refsStr := "[]"
	if err == nil {
		refsStr = string(refsJSON)
	}

	// 4. 执行保存操作（UPSERT）
	_, err = dao.TestTypes.Ctx(ctx).Data(g.Map{
		"id":          id,
		"name":        name,
		"description": desc,
		"references":  refsStr,
		"status":      1,
		"price":       price,
	}).Save()

	if err != nil {
		g.Log().Errorf(ctx, "保存测试类型 [%d] %s 失败: %v", id, name, err)
	} else {
		g.Log().Infof(ctx, "  同步测试类型 [%d] %s (题目: %d, 定价: %.2f)", id, name, questionCount, price)
	}
}

// 辅助函数：将Option数组转换为正确的JSON格式
func optionsToJSON(options []Option) string {
	if len(options) == 0 {
		return "[]"
	}

	// 使用正确的字段顺序：tag, content, score
	var optionList []map[string]interface{}
	for i, opt := range options {
		// 确保每个选项都有正确的标签（A, B, C, D等）
		tag := opt.Tag
		if tag == "" {
			// 如果标签为空，根据索引生成标签
			tag = string(rune('A' + i))
		}
		optionList = append(optionList, map[string]interface{}{
			"tag":     tag,
			"content": opt.Content,
			"score":   opt.Score,
		})
	}

	jsonBytes, err := json.Marshal(optionList)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

// CheckDatabaseStatus 检查数据库状态（可用于健康检查）
func CheckDatabaseStatus() (bool, int, error) {
	ctx := context.Background()

	count, err := dao.TestTypes.Ctx(ctx).Count()
	if err != nil {
		return false, 0, err
	}

	return count > 0, count, nil
}
