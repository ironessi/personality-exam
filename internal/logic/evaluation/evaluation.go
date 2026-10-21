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
	data = normalizeQuestionIDs(data)

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

// normalizeQuestionIDs 把本次作答的题目 ID 折算成从 1 开始的连续序号。
//
// 各个评分函数都是按「题号 1..N」写的：MBTI 写死 case 1,5,9、Personalities
// 写死 case 1..20、RIASEC 用 questionID%6、AMS 用 questionID%2、SBTI 用
// sbtiQuestionMap[1..32]、SCL-90 直接拿 QuestionID 当数组下标。这些写法在
// 题目 ID 恰好从 1 开始的库上成立，但 public.sql 建库后走 /api/initialize，
// 每个测评的题目 ID 都是全局自增的（MBTI 是 376..385，SCL-90 是 1..95），
// 于是 MBTI 四维全 0、Personalities 大五全 0、RIASEC 六个维度整体错位
// （376%6=4，R 变成 S）、AMS 的奇偶被整体翻转（取决于 base 的奇偶）。
//
// 与其在 41 个评分函数里各写一遍偏移，这里统一折算一次：取本次作答的最小
// QuestionID 作基准，映射成 1..N。旧库（从 1 开始）是原样返回，不产生拷贝。
//
// 已知前提：题目 ID 连续、且用户是按题目顺序作答的。缺题会让后续题号前移，
// 这是历史实现就有的假设，不是本次引入的——要修得靠「题目顺序元数据」，
// 属于 SCL-90 题库重写那一批工作。
func normalizeQuestionIDs(data []v1.AnswerRecord) []v1.AnswerRecord {
	if len(data) == 0 {
		return data
	}

	base := data[0].QuestionID
	for _, record := range data[1:] {
		if record.QuestionID < base {
			base = record.QuestionID
		}
	}
	if base == 1 {
		return data
	}

	normalized := make([]v1.AnswerRecord, len(data))
	copy(normalized, data)
	for i := range normalized {
		normalized[i].QuestionID = normalized[i].QuestionID - base + 1
	}
	return normalized
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
