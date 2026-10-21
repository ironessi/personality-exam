package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// Factor 因子定义
type Factor struct {
	Name     string
	Items    []int
	Score    float64
	RawScore int
}

// SCL90TestResult SCL-90测试结果
type SCL90TestResult struct {
	TotalScore      int
	PositiveCount   int
	NegativeCount   int
	Factors         []Factor
	PositiveItems   []int
	HighRiskFactors []Factor
	NeedFurtherEval bool
	Level           string
	Title           string
	Analysis        string
	Suggestions     []string
}

// SCL90Factors SCL-90 量表定义
var SCL90Factors = []Factor{
	{
		Name:  "躯体化",
		Items: []int{1, 4, 12, 27, 40, 42, 48, 49, 52, 53, 56, 58},
	},
	{
		Name:  "强迫症状",
		Items: []int{3, 9, 10, 28, 38, 45, 46, 51, 55, 65},
	},
	{
		Name:  "人际关系敏感",
		Items: []int{6, 21, 34, 36, 37, 41, 61, 69, 73},
	},
	{
		Name:  "抑郁",
		Items: []int{5, 14, 15, 20, 22, 26, 29, 30, 31, 32, 54, 71, 79},
	},
	{
		Name:  "焦虑",
		Items: []int{2, 17, 23, 33, 39, 57, 72, 78, 80, 86},
	},
	{
		Name:  "敌对",
		Items: []int{11, 24, 63, 67, 74, 81},
	},
	{
		Name:  "恐怖",
		Items: []int{13, 25, 47, 50, 70, 75, 82},
	},
	{
		Name:  "偏执",
		Items: []int{8, 18, 43, 68, 76, 83},
	},
	{
		Name:  "精神病性",
		Items: []int{7, 16, 35, 62, 77, 84, 85, 87, 88, 90},
	},
	{
		Name:  "其他",
		Items: []int{19, 44, 59, 60, 64, 66, 89},
	},
}

// SCL90CrisisRules 危机信号题项（1 基题号）与其触发下限。
//
// 为什么要脱离总分与因子均值单独判定：两者都会把单个严重题项稀释掉。
// 「有轻生念头」答“严重”、其余 89 题全答“没有”时，总分只有 94，
// 会被总分分档判成「正常范围」；躯体化因子均值 (11×1+5)/12 ≈ 1.33
// 也够不到 HighRiskFactors 的 2.0 门槛。结果是系统会告诉一个报告了
// 自杀意念的人“你在正常范围内”——这类作答必须无条件升级处理。
var SCL90CrisisRules = map[int]int{
	58: 2, // 有轻生念头：答“很轻”及以上即触发
	20: 4, // 感到生活无意义：答“偏重”及以上触发
	26: 4, // 对未来绝望：答“偏重”及以上触发
	75: 4, // 感到自己毫无价值：答“偏重”及以上触发
}

// SCL90CrisisLabel 危机信号在 high_risk_factors 中的标记名。
const SCL90CrisisLabel = "危机信号（需立即评估）"

// SCL90CrisisLevel 命中危机信号时覆盖 level 的取值。
const SCL90CrisisLevel = "需立即评估"

// getScore 把一次作答折算成 1~5 分。
//
// 前端传的 `answer` 是选项文本，历史上只认 A~E 五个字母；数字字符串与
// 中文频次词（没有/很轻/中等/偏重/严重）都得认，否则同一份量表换个前端
// 就整卷 0 分（实测踩过：answer="5" 被折成 0，总分直接变 0）。
func getScore(answer string) int {
	scoreMap := map[string]int{
		"A": 1, "B": 2, "C": 3, "D": 4, "E": 5,
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5,
		"没有": 1, "很轻": 2, "中等": 3, "偏重": 4, "严重": 5,
	}

	if score, exists := scoreMap[strings.TrimSpace(answer)]; exists {
		return score
	}
	return 0
}

// scoreOf 优先用客户端给的 score 字段，取不到再回退解析 answer。
//
// `score` 是 v1.AnswerRecord 的自带字段，前端算好了就直接用——比再解析一遍
// 文本可靠。只有落在 1~5 之外（比如 0 或没填）时才回退。
func scoreOf(record v1.AnswerRecord) int {
	if record.Score >= 1 && record.Score <= 5 {
		return int(record.Score)
	}
	return getScore(record.Answer)
}

func (e *Evaluation) SCL90(data []v1.AnswerRecord) (interface{}, error) {
	/*
		if len(data) != 90 {
			return nil, errors.New("invalid data")
		}*/

	// 题号是 **1 基**的（normalizeQuestionIDs 已折算成 1..N），必须减 1 再当下标。
	// 原先写 scores[record.QuestionID]：题 1 的答案写进下标 1，全表错位一格，
	// 因子 1 读到的永远是没人写过的下标 0——整卷看似有分，实际每个因子的题
	// 都是隔壁题目的答案。
	scores := make([]int, 90)
	var overflow int
	for _, record := range data {
		if record.QuestionID < 1 {
			continue
		}
		if record.QuestionID > len(scores) {
			overflow++
			continue
		}
		scores[record.QuestionID-1] = scoreOf(record)
	}
	if overflow > 0 {
		// 题库有 95 题、评分只算 90 项（详见 §9 的题库重写项）。多出来的题
		// 不参与计分，但要留下痕迹，免得以后有人以为它们被算进去了。
		g.Log().Warningf(context.Background(),
			"SCL-90 有 %d 道题的题号超出 90 项计分范围，已忽略", overflow)
	}

	result := calculateDetailedResults(scores)

	// 转换为统一的返回格式
	factorDetails := make(map[string]interface{})
	for _, factor := range result.Factors {
		factorDetails[factor.Name] = map[string]interface{}{
			"score":     factor.Score,
			"raw_score": factor.RawScore,
			"items":     factor.Items,
		}
	}

	highRiskFactors := make([]string, 0)
	for _, factor := range result.HighRiskFactors {
		highRiskFactors = append(highRiskFactors, factor.Name)
	}

	return map[string]interface{}{
		"test_name":         "SCL-90症状自评量表",
		"total_score":       result.TotalScore,
		"positive_count":    result.PositiveCount,
		"negative_count":    result.NegativeCount,
		"max_score":         450, // 90题 * 5分
		"level":             result.Level,
		"title":             result.Title,
		"analysis":          result.Analysis,
		"suggestions":       result.Suggestions,
		"factors":           factorDetails,
		"high_risk_factors": highRiskFactors,
		"positive_items":    result.PositiveItems,
		"need_further_eval": result.NeedFurtherEval,
	}, nil
}

// 详细的计算结果
func calculateDetailedResults(scores []int) SCL90TestResult {
	var result SCL90TestResult

	// 计算总分和阳性项目数
	for i, score := range scores {
		result.TotalScore += score
		if score >= 2 {
			result.PositiveCount++
			result.PositiveItems = append(result.PositiveItems, i+1)
		} else {
			result.NegativeCount++
		}
	}

	// 计算各因子分
	for i := range SCL90Factors {
		factor := &SCL90Factors[i]
		var factorTotal int

		for _, itemNum := range factor.Items {
			factorTotal += scores[itemNum-1] // 项目编号从1开始，数组索引从0开始
		}

		factor.RawScore = factorTotal
		factor.Score = float64(factorTotal) / float64(len(factor.Items))

		// 记录高风险因子（因子分≥2）
		if factor.Score >= 2.0 {
			result.HighRiskFactors = append(result.HighRiskFactors, *factor)
		}
	}

	result.Factors = SCL90Factors

	// 判断是否需要进一步评估：总分≥160 或 任一因子分≥2 或 阳性项目数>43
	result.NeedFurtherEval = result.TotalScore >= 160 || len(result.HighRiskFactors) > 0 || result.PositiveCount > 43

	// 根据总分进行分级。分界点用 < 160，与上面 NeedFurtherEval 的 ">= 160"
	// 严格互补：原先写 <= 160，于是总分正好 160 时会一边报「正常范围」、
	// 一边 need_further_eval=true，同一份结果自相矛盾。
	if result.TotalScore < 160 {
		result.Level = "正常范围"
		result.Title = "情绪状态正常"
		result.Analysis = "您的SCL-90总分为正常范围，阳性项目数在可控范围内。"
		result.Suggestions = []string{
			"继续保持健康的生活方式",
			"定期关注心理健康",
			"如有不适及时寻求帮助",
		}
	} else if result.TotalScore <= 225 {
		result.Level = "轻度异常"
		result.Title = "建议关注心理健康"
		result.Analysis = "您的SCL-90总分显示轻度异常，建议关注心理健康。"
		result.Suggestions = []string{
			"学习压力管理技巧",
			"适当减少压力源",
			"考虑心理咨询",
			"建立健康的生活习惯",
		}
	} else if result.TotalScore <= 315 {
		result.Level = "明显异常"
		result.Title = "建议寻求专业评估"
		result.Analysis = "您的SCL-90总分显示明显异常，建议寻求专业评估。"
		result.Suggestions = []string{
			"寻求专业心理咨询师评估",
			"结合临床访谈综合判断",
			"建立支持系统",
			"必要时考虑心理治疗",
		}
	} else {
		result.Level = "严重异常"
		result.Title = "强烈建议专业帮助"
		result.Analysis = "您的SCL-90总分显示严重异常，强烈建议立即寻求专业帮助。"
		result.Suggestions = []string{
			"立即寻求精神科医生帮助",
			"考虑药物治疗配合心理治疗",
			"建立紧急支持网络",
			"学习危机应对策略",
		}
	}

	// 总分不高、但个别因子超过 2 分（或阳性项目数偏多）时，不能继续报「正常范围」：
	// 此时 need_further_eval 已为 true，两份结论会打架。这里只做「不再说正常」的
	// 升级，不新增更重的档位——因子分是描述性指标，靠它推严重程度是过度解读。
	if result.Level == "正常范围" && result.NeedFurtherEval {
		result.Level = "轻度异常"
		result.Title = "总分不高，但有个别方面需要关注"
		if n := len(result.HighRiskFactors); n > 0 {
			result.Analysis = fmt.Sprintf(
				"您的SCL-90总分未达到 160 分的筛查阳性线，但有 %d 个因子分达到或超过 2 分。"+
					"因子分高只说明这一组体验在近期比较突出，并不等于患病；"+
					"建议结合最近的压力事件和生活状态观察一段时间，若持续影响工作、睡眠或人际，"+
					"再与精神科医生或注册心理师面谈评估。", n)
		} else {
			result.Analysis = fmt.Sprintf(
				"您的SCL-90总分未达到 160 分的筛查阳性线，但阳性项目数（%d 项）偏多，"+
					"说明有不少条目报告了轻度以上的不适。单项都很轻时总分不会高，"+
					"但分布广本身值得留意；建议留意近期压力来源，必要时寻求专业评估。",
				result.PositiveCount)
		}
		result.Suggestions = []string{
			"留意近期压力事件与睡眠、作息的变化",
			"把困扰记录下来，观察两周是否有缓解",
			"若持续影响工作或人际，考虑预约心理咨询",
		}
	}

	// 如果有高风险因子，添加额外的分析
	if len(result.HighRiskFactors) > 0 {
		result.Analysis += " 检测到以下高风险因子："
		for _, factor := range result.HighRiskFactors {
			result.Analysis += factor.Name + "(" + fmt.Sprintf("%.2f", factor.Score) + "分)、"
		}
		result.Analysis = result.Analysis[:len(result.Analysis)-1] + "，需要特别关注。"
	}

	// 危机信号单独判定，不受总分与因子均值稀释（见 SCL90CrisisRules 注释）。
	// 这一段必须放在总分分档之后，因为它的结论要覆盖分档结论。
	var crisisItems []int
	for itemNum, floor := range SCL90CrisisRules {
		if scores[itemNum-1] >= floor {
			crisisItems = append(crisisItems, itemNum)
		}
	}
	sort.Ints(crisisItems)

	if len(crisisItems) > 0 {
		result.NeedFurtherEval = true
		result.HighRiskFactors = append(result.HighRiskFactors, Factor{Name: SCL90CrisisLabel})
		result.Level = SCL90CrisisLevel
		result.Title = "检测到危机信号，请立即寻求专业帮助"
		result.Analysis = fmt.Sprintf(
			"本次作答中有 %d 个题项达到危机信号标准（题号 %v）。"+
				"这类情况不适合继续自行观察，请尽快联系精神科医生或注册心理师面谈评估；"+
				"若出现自伤或自杀的念头、计划或行为，请立即联系当地急救与心理援助资源，"+
				"或前往最近的正规医疗机构急诊。",
			len(crisisItems), crisisItems)
		result.Suggestions = append([]string{
			"立即寻求专业帮助，不要等待",
			"联系当地急救与心理援助资源，或前往最近的正规医疗机构急诊",
			"让身边可信任的人知晓你的情况，尽量不要独处",
		}, result.Suggestions...)
	}

	return result
}
