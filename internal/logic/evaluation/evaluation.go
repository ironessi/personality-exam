package evaluation

import (
	v1 "PersonalityTest/api/test_type/v1"
	"PersonalityTest/internal/consts"
)

type Evaluation struct{}

func New() *Evaluation {
	return &Evaluation{}
}

// 统一的评估入口
func (e *Evaluation) Evaluate(testType int, data []v1.AnswerRecord) (interface{}, error) {
	// 首先检查增强的分析函数
	if analyzeFunc, exists := EnhancedAnalysisFunctions[testType]; exists {
		return analyzeFunc(data)
	}

	// 使用原有的分析函数
	switch testType {
	case consts.SCL90:
		return e.SCL90(data)
	case consts.GeneralLifeSatisfaction:
		return e.GeneralLifeSatisfaction(data)
	case consts.PerceivedStress:
		return e.PerceivedStress(data)
	case consts.PROMISAnxiety:
		return e.PROMISAnxiety(data)
	case consts.SocialDistress:
		return e.SocialDistress(data)
	case consts.UCLALoneliness:
		return e.UCLALoneliness(data)
	case consts.LifeMeaning:
		return e.LifeMeaning(data)
	case consts.Panas:
		return e.Panas(data) // 注意：这里应该是 Panas，不是 PANAS
	case consts.SelfEfficacy:
		return e.SelfEfficacy(data)
	case consts.PhysicalAggression:
		return e.PhysicalAggression(data)
	case consts.Hostility:
		return e.Hostility(data)
	// 职业测试
	case consts.RIASEC:
		return e.RIASEC(data)
	case consts.MBTI:
		return e.MBTI(data)
	case consts.Personalities:
		return e.Personalities(data)
	case consts.PSTRI:
		return e.PSTRI(data)
	case consts.PhilPersonality:
		return e.PhilPersonality(data)
	// 趣味测试
	case consts.STScale:
		return e.SocialStatusIndex(data) // 复用社会地位指数测试的评分逻辑
	case consts.AMS:
		return e.AMS(data)
	case consts.Zootopia:
		return e.Zootopia(data)
	case consts.StanfordBinet:
		return e.StanfordBinet(data)
	case consts.Rorschach:
		return e.Rorschach(data)
	// 新增心理测评问卷
	case consts.BDI:
		return e.BDI(data)
	case consts.SDSAnxiety:
		return e.SDSAnxiety(data)
	case consts.HCL32:
		return e.HCL32(data)
	case consts.MHLAQ:
		return e.MHLAQ(data)
	case consts.TwentyNineDegreePersonality:
		return e.TwentyNineDegreePersonality(data)
	case consts.EMP:
		return e.EMP(data)
	case consts.DISC:
		return e.DISC(data)
	case consts.EPQ:
		return e.EPQ(data)
	case consts.SocialStatusIndex:
		return e.SocialStatusIndex(data)
	case consts.Machiavellianism:
		return e.Machiavellianism(data)
	// 新增心理测评问卷
	case consts.SelfControlIndex:
		return e.SelfControlIndex(data)
	case consts.PassivePersonality:
		return e.PassivePersonality(data)
	case consts.TrueColors:
		return e.TrueColors(data)
	case consts.ObedienceTest:
		return e.ObedienceTest(data)
	case consts.MMPI:
		return e.MMPI(data)
	case consts.AnxietyScale:
		return e.AnxietyScale(data)
	case consts.EPPS:
		return e.EPPS(data)
	case consts.HSP:
		return e.HSP(data)
	case consts.NPI:
		return e.NPI(data)
	case consts.BemSexRoleInventory:
		return e.BemSexRoleInventory(data)
	case consts.SBTI:
		return e.SBTI(data)
	default:
		// 默认分析
		return e.defaultAnalysis(data)
	}
}

// 默认分析函数
func (e *Evaluation) defaultAnalysis(data []v1.AnswerRecord) (interface{}, error) {
	totalScore := 0
	for _, record := range data {
		totalScore += int(record.Score)
	}

	return map[string]interface{}{
		"test_name":   "心理测评测试",
		"total_score": totalScore,
		"level":       "正常",
		"title":       "测试完成",
		"analysis":    "测试分析完成，请查看详细结果。",
		"suggestions": []string{
			"保持积极的生活态度",
			"定期进行心理健康评估",
			"寻求专业心理咨询帮助",
		},
	}, nil
}

// 如果需要，可以添加一个映射表用于增强分析函数
var EnhancedAnalysisFunctions = map[int]func([]v1.AnswerRecord) (interface{}, error){
	// 这里可以注册增强的分析函数
}
