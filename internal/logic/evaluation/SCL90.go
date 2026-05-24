package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"fmt"
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

func getScore(answer string) int {
	scoreMap := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
	}

	if score, exists := scoreMap[answer]; exists {
		return score
	}
	return 0
}

func (e *Evaluation) SCL90(data []v1.AnswerRecord) (interface{}, error) {
	/*
		if len(data) != 90 {
			return nil, errors.New("invalid data")
		}*/

	scores := make([]int, 90)
	for _, record := range data {
		if record.QuestionID >= 0 && record.QuestionID < 90 {
			scores[record.QuestionID] = getScore(record.Answer)
		}
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

	// 根据总分进行分级
	if result.TotalScore <= 160 {
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

	// 如果有高风险因子，添加额外的分析
	if len(result.HighRiskFactors) > 0 {
		result.Analysis += " 检测到以下高风险因子："
		for _, factor := range result.HighRiskFactors {
			result.Analysis += factor.Name + "(" + fmt.Sprintf("%.2f", factor.Score) + "分)、"
		}
		result.Analysis = result.Analysis[:len(result.Analysis)-1] + "，需要特别关注。"
	}

	return result
}
