package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
	"math"
)

// SBTITemplate 定义了人格模板的结构
type SBTITemplate struct {
	Name     string
	Type     string
	Levels   []string // 15维特征向量，例如 ["L", "H", "M", ...]
	Analysis string
}

// 预定义部分典型 SBTI 人格模板作为示例（完整需配置25种）
var sbtiTemplates = []SBTITemplate{
	{
		Name: "吗喽 (MALOU)",
		Type: "MALOU",
		Levels: []string{"L", "L", "L", "M", "H", "H", "H", "M", "M", "M", "L", "L", "H", "H", "L"},
		Analysis: "高度自嘲、适应力极强，表面上是个“吗喽”，内心其实向往自由与平静，执行力相对佛系。",
	},
	{
		Name: "伪人 (FAKE)",
		Type: "FAKE",
		Levels: []string{"M", "L", "H", "L", "L", "M", "H", "H", "M", "H", "M", "H", "L", "H", "L"},
		Analysis: "戴着厚重的社交面具，人前人后反差极大。内心有极强的边界感，对外人保持防备。",
	},
	{
		Name: "掌控者 (BOSS)",
		Type: "BOSS",
		Levels: []string{"H", "H", "H", "M", "M", "L", "L", "H", "M", "H", "H", "H", "H", "L", "H"},
		Analysis: "执行力爆表，有着强烈的核心价值驱动，不喜欢任何计划外的变数。",
	},
	{
		Name: "送钱者 (ATM-er)",
		Type: "ATM-er",
		Levels: []string{"L", "M", "M", "H", "H", "H", "M", "M", "L", "L", "L", "M", "H", "L", "M"},
		Analysis: "高情感投入，极易变成恋爱脑或付出型人格，个人边界感较弱。",
	},
	{
		Name: "死皮赖脸 (DEAD)",
		Type: "DEAD",
		Levels: []string{"L", "L", "L", "L", "L", "L", "L", "L", "L", "L", "L", "L", "L", "L", "L"},
		Analysis: "极度虚无主义，对什么都提不起兴趣，生活状态如同“死掉”一般。",
	},
	{
		Name: "小丑 (CLOWN)",
		Type: "CLOWN",
		Levels: []string{"L", "M", "M", "L", "H", "H", "H", "L", "M", "M", "M", "H", "L", "M", "H"},
		Analysis: "总是习惯用幽默掩饰内心的伤痕，人前活跃气氛，人后独自舔舐伤口。",
	},
	{
		Name: "路人甲 (NPC)",
		Type: "NPC",
		Levels: []string{"M", "M", "M", "M", "M", "M", "M", "M", "M", "M", "M", "M", "M", "M", "M"},
		Analysis: "完美的中庸之道，不冒进、不退缩，默默无闻地完成生活布置的每一个微小任务。",
	},
	{
		Name: "幽灵 (GHOST)",
		Type: "GHOST",
		Levels: []string{"M", "L", "L", "L", "L", "L", "M", "M", "H", "H", "L", "M", "L", "L", "L"},
		Analysis: "存在感极低，经常被人遗忘，自己也乐于躲在角落里偷偷观察世界。",
	},
	{
		Name: "小太阳 (SUN)",
		Type: "SUN",
		Levels: []string{"H", "H", "H", "H", "H", "H", "H", "H", "H", "H", "H", "H", "H", "H", "L"},
		Analysis: "充满正能量的化身，热衷于温暖他人，但有时候会忽略自己内心的阴暗面。",
	},
	{
		Name: "月亮翁 (MOON)",
		Type: "MOON",
		Levels: []string{"L", "L", "M", "H", "M", "H", "L", "L", "L", "M", "L", "H", "L", "L", "H"},
		Analysis: "敏感、感性、思维跳跃。通常晚睡，情绪起伏像潮汐一样有规律但难以捉摸。",
	},
	{
		Name: "修勾 (DOGE)",
		Type: "DOGE",
		Levels: []string{"M", "H", "H", "L", "H", "H", "H", "M", "M", "M", "L", "H", "H", "L", "M"},
		Analysis: "忠诚度极高，对感情执着，渴望被摸摸头，一旦认定就会死心塌地。",
	},
	{
		Name: "高冷猫 (CAT)",
		Type: "CAT",
		Levels: []string{"H", "H", "L", "H", "L", "L", "L", "H", "H", "H", "H", "M", "H", "L", "H"},
		Analysis: "极大概率是社恐，边界感极强，偶尔需要陪伴，但绝不能过度打扰。",
	},
	{
		Name: "打工人 (WORKER)",
		Type: "WORKER",
		Levels: []string{"L", "L", "M", "M", "H", "L", "L", "M", "M", "M", "M", "H", "M", "M", "L"},
		Analysis: "为了几两碎银子折腰，常常抱怨但又总是准时打卡，被生活捶打得没有脾气。",
	},
	{
		Name: "睡神 (SLEEPER)",
		Type: "SLEEPER",
		Levels: []string{"L", "L", "L", "L", "H", "L", "L", "L", "M", "L", "L", "M", "L", "L", "L"},
		Analysis: "万事皆可睡，遇到压力和困难的第一反应就是闭上眼睛逃避。",
	},
	{
		Name: "干饭人 (EATER)",
		Type: "EATER",
		Levels: []string{"M", "M", "L", "M", "H", "M", "L", "L", "M", "M", "M", "H", "L", "H", "M"},
		Analysis: "吃是生活中最大的动力，信奉“没有什么是一顿火锅解决不了的”。",
	},
	{
		Name: "酒门信徒 (DRINKER)",
		Type: "DRINKER",
		Levels: []string{"L", "L", "L", "H", "H", "L", "H", "L", "M", "L", "L", "H", "M", "M", "M"},
		Analysis: "酒精沙场老将，越喝越清醒，只在微醺状态下吐露真心。",
	},
	{
		Name: "国王 (KING)",
		Type: "KING",
		Levels: []string{"H", "H", "H", "H", "M", "L", "H", "H", "H", "H", "H", "H", "M", "H", "M"},
		Analysis: "天生的领导者，要求高，掌控欲强，难以容忍下放权力。",
	},
	{
		Name: "弄臣 (JESTER)",
		Type: "JESTER",
		Levels: []string{"M", "H", "L", "H", "M", "M", "L", "L", "L", "M", "H", "L", "M", "L", "M"},
		Analysis: "虽然看似不羁，其实深谙世故，懂得如何在复杂的人际网中明哲保身。",
	},
	{
		Name: "浪漫诗人 (POET)",
		Type: "POET",
		Levels: []string{"L", "L", "L", "H", "H", "H", "H", "L", "M", "M", "L", "H", "L", "L", "L"},
		Analysis: "为了浪漫可以放弃一切，理想主义者，现实生活中的重度拖延症患者。",
	},
	{
		Name: "机器人 (ROBOT)",
		Type: "ROBOT",
		Levels: []string{"H", "H", "H", "L", "L", "L", "H", "H", "H", "H", "H", "H", "H", "L", "M"},
		Analysis: "无情的效率机器，把感情看作是影响运算速度的bug。",
	},
	{
		Name: "守护骑士 (KNIGHT)",
		Type: "KNIGHT",
		Levels: []string{"H", "H", "L", "H", "H", "H", "H", "H", "L", "L", "M", "H", "H", "M", "M"},
		Analysis: "极具正义感，遇到不公必会发声，哪怕有时候会吃亏。",
	},
	{
		Name: "巨婴 (BABY)",
		Type: "BABY",
		Levels: []string{"L", "L", "L", "H", "H", "H", "H", "H", "L", "L", "L", "L", "L", "H", "H"},
		Analysis: "需要被无限包容与照顾，遇到困难总希望有人能帮自己扛。",
	},
	{
		Name: "外星人 (ALIEN)",
		Type: "ALIEN",
		Levels: []string{"L", "H", "L", "H", "M", "L", "H", "L", "L", "L", "M", "H", "M", "L", "M"},
		Analysis: "脑回路清奇，在这个地球上常常感到格格不入。",
	},
	{
		Name: "忍者 (NINJA)",
		Type: "NINJA",
		Levels: []string{"M", "H", "L", "L", "M", "L", "L", "H", "M", "M", "M", "H", "L", "H", "L"},
		Analysis: "极能隐忍，喜怒不形于色，可以在极端的环境下生存。",
	},
	{
		Name: "活菩萨 (GOD)",
		Type: "GOD",
		Levels: []string{"H", "L", "H", "H", "H", "H", "H", "L", "L", "L", "H", "H", "H", "H", "L"},
		Analysis: "极其善良，几乎不会拒绝别人，常常也是被欺负的那个。",
	},
}

// 题目到维度的映射（由于前端题目顺序可能打乱，且无法控制题号，因此根据我们解析的 JSON 设定一个简化的映射库）
// 假设前端传入的 QuestionID 为 1-32，我们将其映射至这 15 个维度（包含隐藏门控）
var sbtiQuestionMap = map[int]string{
	1: "S1", 2: "S1",
	3: "S2", 4: "S2",
	5: "S3", 6: "S3",
	7: "E1", 8: "E1",
	9: "E2", 10: "E2",
	11: "E3", 12: "E3",
	13: "A1", 14: "A1",
	15: "A2", 16: "A2",
	17: "A3", 18: "A3",
	19: "Ac1", 20: "Ac1",
	21: "Ac2", 22: "Ac2",
	23: "Ac3", 24: "Ac3",
	25: "So1", 26: "So1",
	27: "So2", 28: "So2",
	29: "So3", 30: "So3",
	31: "Gate", // 爱好选择（门控）
	32: "Drunk", // 饮酒态度触发
}

// 维度顺序，用于向量对比
var sbtiDimensions = []string{
	"S1", "S2", "S3",
	"E1", "E2", "E3",
	"A1", "A2", "A3",
	"Ac1", "Ac2", "Ac3",
	"So1", "So2", "So3",
}

// SBTI 傻大个性格测试评估逻辑
func (e *Evaluation) SBTI(data []v1.AnswerRecord) (interface{}, error) {
	// 用于存放 15 个维度的总分
	scores := make(map[string]int)
	for _, dim := range sbtiDimensions {
		scores[dim] = 0 // 初始化为 0
	}

	totalScore := 0
	isDrunk := false

	// 遍历所有用户答题，并累计分数
	for _, record := range data {
		qid := record.QuestionID
		score := int(record.Score)
		totalScore += score

		// 获取该题对应的维度
		if dim, exists := sbtiQuestionMap[qid]; exists {
			// 隐藏门控判定 (第32题：如果选了B(即score为2)——白酒灌保温杯，则触发隐藏人格)
			if dim == "Drunk" && score == 2 {
				isDrunk = true
			}

			// 普通维度累计分值
			if _, ok := scores[dim]; ok {
				scores[dim] += score
			}
		}
	}

	// 如果触发了 DRUNK 隐藏机制，直接返回特定结果，终止正常计算
	if isDrunk {
		return map[string]interface{}{
			"test_name":   "SBTI 傻大个性格测试",
			"total_score": totalScore,
			"level":       "隐藏触发",
			"title":       "酒精成精 (DRUNK)",
			"analysis":    "您触发了特殊隐藏人格——DRUNK！酒精是您忠实的伴侣，白酒装保温杯是您的信仰。常规的心理维度已经无法定义您...",
			"suggestions": []string{
				"建议适量饮酒，保重身体",
				"别再用白酒保温杯伪装喝水了",
			},
			"vector": []string{},
		}, nil
	}

	// 1. 各维度归一化 (Leveing)：将分数映射为 L, M, H 三个档位
	// （由于每维度一般2题，每题1-3分，维度得分范围为 2-6分）
	userVector := make([]string, len(sbtiDimensions))
	for i, dim := range sbtiDimensions {
		s := scores[dim]
		var level string
		if s <= 3 {
			level = "L"
		} else if s == 4 {
			level = "M"
		} else {
			level = "H"
		}
		userVector[i] = level
	}

	// 2. 匹配最近的模板库 (最小曼哈顿距离)
	minDistance := math.MaxInt32
	var bestMatch SBTITemplate

	for _, tpl := range sbtiTemplates {
		dist := calculateSBTIDistance(userVector, tpl.Levels)
		if dist < minDistance {
			minDistance = dist
			bestMatch = tpl
		}
	}

	analysisContent := fmt.Sprintf("经15维度矩阵测算，您的基础心智模型最接近 【%s】。\n\n【核心特质解析】\n%s\n\n【十五特征参量】\n%v\n与该模板相似度偏差值：%d",
		bestMatch.Name, bestMatch.Analysis, userVector, minDistance)

	return map[string]interface{}{
		"test_name":   "SBTI 傻大个性格测试",
		"total_score": totalScore,
		"level":       bestMatch.Type,
		"title":       bestMatch.Name,
		"analysis":    analysisContent,
		"suggestions": []string{
			"本测试仅供娱乐调侃，不可作为严谨的心理学诊断参考",
			"不必对结果过度对号入座，生活仍充满变数",
			"将结果发朋友圈，看有多少同款“吗喽”或是“伪人”！",
		},
		"vector": userVector,
	}, nil
}

// calculateSBTIDistance 运用类似于曼哈顿距离的算法计算相似度
func calculateSBTIDistance(uVec, tVec []string) int {
	dist := 0
	valMap := map[string]int{"L": 1, "M": 2, "H": 3}

	for i := 0; i < len(uVec) && i < len(tVec); i++ {
		uVal := valMap[uVec[i]]
		tVal := valMap[tVec[i]]

		diff := uVal - tVal
		if diff < 0 {
			diff = -diff
		}
		dist += diff
	}
	return dist
}
