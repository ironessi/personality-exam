package consts

import "sort"

// 本文件是「测评元数据注册表」——42 个测评的**唯一**元数据来源（ADR-0008）。
//
// 为什么要有它：在此之前「这个测评该检索哪部分知识库」这件事根本无处可查——
// manifests/config/config.yaml 里只有一个手写的灰度白名单 testTypes: [13, 1, 28]，
// 而 42 个评分文件内部散落着 103 处写死的阈值。新增一个测评要同时改语料、改检索、
// 改契约，于是「覆盖全部 42 个测评」变成一件需要 42 次重复劳动的事。
//
// 现在的形状：测评 → {构念域, 增强策略}，检索范围由本表算好后放进请求体的
// test.domains，Python 侧**不判断测评类型**，只照 domains 过滤。
//
// ⚠️ 阈值（分档边界 / 量表上界）**还没有**搬进本文件，见文末 TODO。

// ── 知识库分区键：构念域（ADR-0008）─────────────────────────────
//
// namespace **不是测评名**。按测评建 namespace 会得到 42 个几乎同质的索引，
// 而 chunk_id 的前缀就是 namespace——两套命名混用时 Chroma 的 upsert 会相互
// 覆盖且**不报错**（GLOSSARY 第 4 条歧义）。
//
// 构念域是**封闭集合**：新增测评只需在 TestRegistry 里登记它落在哪个域，
// 语料目录、检索代码、Python 契约都不用改。

// Domain 是知识库分区键，取值见下面 8 个常量。
type Domain string

const (
	DomainSymptom   Domain = "symptom"   // 症状严重度：SCL-90 / BDI / SDS / HCL-32 / MMPI …
	DomainTrait     Domain = "trait"     // 特质人格：EPQ / NPI / HSP / EPPS / 马基雅维利主义 …
	DomainTypology  Domain = "typology"  // 类型学人格：MBTI / DISC / True Colors / 性别角色
	DomainSocial    Domain = "social"    // 人际与社会：孤独感 / 社交困扰 / 同理心 / 服从性
	DomainWellbeing Domain = "wellbeing" // 福祉与积极心理：生活满意度 / 生命意义 / 情感平衡 / 自我效能
	DomainCareer    Domain = "career"    // 职业兴趣：RIASEC
	DomainAbility   Domain = "ability"   // 能力与智力：心理健康素养（本域语料待写，P1）
	DomainCommon    Domain = "common"    // 通用层：任何测评都可检索
)

// AllDomains 是全部合法构念域。顺序与 data/raw 子目录名、
// rag-service/app/core/config.py 的 kb_namespaces 默认值保持一致。
var AllDomains = []Domain{
	DomainSymptom,
	DomainTrait,
	DomainTypology,
	DomainSocial,
	DomainWellbeing,
	DomainCareer,
	DomainAbility,
	DomainCommon,
}

// IsValidDomain 供注册表自检与外部输入校验使用。
func IsValidDomain(d Domain) bool {
	for _, v := range AllDomains {
		if v == d {
			return true
		}
	}
	return false
}

// ── 增强策略（ADR-0008）─────────────────────────────────────────
//
// ⚠️ 这不是性能开关。它存在的唯一理由是：
// **不要对不该解释的东西给出看起来很专业的解释。**

// EnhancePolicy 决定一个测评接入检索增强的方式。
type EnhancePolicy string

const (
	// PolicyFull = 检索「本域 + common」，正常生成深度解读。38 个测评走这条。
	PolicyFull EnhancePolicy = "full"

	// PolicyCommonOnly = 只检索 common，且强制带上 needs_referral。
	//
	// 用于两个测验：斯坦福-比奈智力量表、罗夏墨迹测验。它们的现行版本由出版方
	// 持有版权，标准化施测与计分（如罗夏的 Exner 综合系统）**必须由受训主试
	// 一对一实施**；本平台的实现是自评近似（见 internal/logic/evaluation/Rorschach.go
	// 里用总分区间判「外向情感型」）。对它生成个性化深度解读等于**伪造权威**。
	PolicyCommonOnly EnhancePolicy = "common_only"

	// PolicyDisabled = 不接入增强，提交后 result_data 里不会出现 rag 字段。
	//
	// 用于两个娱乐测评：疯狂动物城性格测试、SBTI。给它们写深度心理学解读既误导
	// 读者，又涉及迪士尼 IP。
	PolicyDisabled EnhancePolicy = "disabled"
)

// ── 注册表 ──────────────────────────────────────────────────────

// TestMeta 是一个测评的全部增强元数据。
//
// Key / Name 只是**标识与展示**，都不是分区键；分区键是 Domain。
type TestMeta struct {
	Code   int           // consts.go 里的数字编码
	Key    string        // 机读名，如 "mbti"；用于日志与缓存键，**不是 namespace**
	Name   string        // 展示名；⚠️ 仅在拿不到 DB test_types.name 时兜底
	Domain Domain        // 构念域；策略为 disabled 时为空
	Policy EnhancePolicy // 见上
}

// TestRegistry 是 42 个测评的权威映射（附录 B 的唯一副本）。
//
// 改动本表时必须同步三处：doc/RAG增强系统开发文档.md 附录 B、
// data/raw/<域>/SOURCES.md 的覆盖清单、rag-service 侧无（它不判断测评类型）。
var TestRegistry = map[int]TestMeta{
	// ── 症状严重度（11）─────────────────────────────────────────
	SCL90: {
		Code: SCL90, Key: "scl90", Name: "SCL-90 症状自评量表",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	PerceivedStress: {
		Code: PerceivedStress, Key: "perceived_stress", Name: "压力知觉量表",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	PROMISAnxiety: {
		Code: PROMISAnxiety, Key: "promis_anxiety", Name: "PROMIS 焦虑量表",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	PhysicalAggression: {
		Code: PhysicalAggression, Key: "physical_aggression", Name: "身体攻击性测试",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	Hostility: {
		Code: Hostility, Key: "hostility", Name: "敌意测试",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	PSTRI: {
		Code: PSTRI, Key: "pstri", Name: "PSTRI 压力知觉测试量表",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	BDI: {
		Code: BDI, Key: "bdi", Name: "贝克抑郁自评量表",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	SDSAnxiety: {
		Code: SDSAnxiety, Key: "sds_anxiety", Name: "焦虑症自测量表（SDS）",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	HCL32: {
		Code: HCL32, Key: "hcl32", Name: "轻躁狂症状自评量表（HCL-32）",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	MMPI: {
		Code: MMPI, Key: "mmpi", Name: "多相人格测试（MMPI）",
		Domain: DomainSymptom, Policy: PolicyFull,
	},
	AnxietyScale: {
		Code: AnxietyScale, Key: "anxiety_scale", Name: "焦虑性测试（As 量表）",
		Domain: DomainSymptom, Policy: PolicyFull,
	},

	// ── 特质人格（11）───────────────────────────────────────────
	Personalities: {
		Code: Personalities, Key: "personalities", Name: "Personalities 人格测试",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	PhilPersonality: {
		Code: PhilPersonality, Key: "phil_personality", Name: "菲尔人格十项测试",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	AMS: {
		Code: AMS, Key: "ams", Name: "成就动机测试（AMS）",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	TwentyNineDegreePersonality: {
		Code: TwentyNineDegreePersonality, Key: "twenty_nine_degree", Name: "29度人格测试",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	EPQ: {
		Code: EPQ, Key: "epq", Name: "艾森克人格问卷（EPQ）",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	Machiavellianism: {
		Code: Machiavellianism, Key: "machiavellianism", Name: "马基雅维利人格测试",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	SelfControlIndex: {
		Code: SelfControlIndex, Key: "self_control_index", Name: "自我控力指数评估",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	PassivePersonality: {
		Code: PassivePersonality, Key: "passive_personality", Name: "被动型人格测试",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	EPPS: {
		Code: EPPS, Key: "epps", Name: "爱德华个人偏好测验（EPPS）",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	HSP: {
		Code: HSP, Key: "hsp", Name: "敏感度测试（HSP）",
		Domain: DomainTrait, Policy: PolicyFull,
	},
	NPI: {
		Code: NPI, Key: "npi", Name: "自恋人格量表（NPI）",
		Domain: DomainTrait, Policy: PolicyFull,
	},

	// ── 类型学人格（4）──────────────────────────────────────────
	MBTI: {
		Code: MBTI, Key: "mbti", Name: "16型人格测试（MBTI）",
		Domain: DomainTypology, Policy: PolicyFull,
	},
	DISC: {
		Code: DISC, Key: "disc", Name: "DISC 性格测试",
		Domain: DomainTypology, Policy: PolicyFull,
	},
	TrueColors: {
		Code: TrueColors, Key: "true_colors", Name: "色彩性格测试（True Colors）",
		Domain: DomainTypology, Policy: PolicyFull,
	},
	BemSexRoleInventory: {
		Code: BemSexRoleInventory, Key: "bem_sex_role", Name: "贝姆性别角色测试量表",
		Domain: DomainTypology, Policy: PolicyFull,
	},

	// ── 人际与社会（6）──────────────────────────────────────────
	SocialDistress: {
		Code: SocialDistress, Key: "social_distress", Name: "社交困扰量表",
		Domain: DomainSocial, Policy: PolicyFull,
	},
	UCLALoneliness: {
		Code: UCLALoneliness, Key: "ucla_loneliness", Name: "UCLA 孤独感量表",
		Domain: DomainSocial, Policy: PolicyFull,
	},
	STScale: {
		Code: STScale, Key: "st_scale", Name: "社会地位指数测试（ST 量表）",
		Domain: DomainSocial, Policy: PolicyFull,
	},
	EMP: {
		Code: EMP, Key: "emp", Name: "同理心通情量表",
		Domain: DomainSocial, Policy: PolicyFull,
	},
	SocialStatusIndex: {
		// 注意：STScale(17) 的评分逻辑复用本测评（evaluation.go 的 switch 分支），
		// 但两者是**两个测评、两条注册表记录**，不能合并。
		Code: SocialStatusIndex, Key: "social_status_index", Name: "社会地位指数测试",
		Domain: DomainSocial, Policy: PolicyFull,
	},
	ObedienceTest: {
		Code: ObedienceTest, Key: "obedience", Name: "服从性测试",
		Domain: DomainSocial, Policy: PolicyFull,
	},

	// ── 福祉与积极心理（4）──────────────────────────────────────
	GeneralLifeSatisfaction: {
		Code: GeneralLifeSatisfaction, Key: "life_satisfaction", Name: "总体生活满意度量表",
		Domain: DomainWellbeing, Policy: PolicyFull,
	},
	LifeMeaning: {
		Code: LifeMeaning, Key: "life_meaning", Name: "生命意义感测试",
		Domain: DomainWellbeing, Policy: PolicyFull,
	},
	Panas: {
		Code: Panas, Key: "panas", Name: "积极情感消极情感量表",
		Domain: DomainWellbeing, Policy: PolicyFull,
	},
	SelfEfficacy: {
		Code: SelfEfficacy, Key: "self_efficacy", Name: "自我效能感测试",
		Domain: DomainWellbeing, Policy: PolicyFull,
	},

	// ── 职业兴趣（1）────────────────────────────────────────────
	RIASEC: {
		Code: RIASEC, Key: "riasec", Name: "霍兰德职业兴趣测试",
		Domain: DomainCareer, Policy: PolicyFull,
	},

	// ── 能力与智力（1 个走 full）────────────────────────────────
	MHLAQ: {
		Code: MHLAQ, Key: "mhlaq", Name: "心理健康素养量表",
		Domain: DomainAbility, Policy: PolicyFull,
	},

	// ── 只检索通用层（2）────────────────────────────────────────
	StanfordBinet: {
		// 登记在 ability 域只是为了「它属于哪个构念」有据可查；检索时被 common 取代。
		Code: StanfordBinet, Key: "stanford_binet", Name: "斯坦福-比奈智力量表",
		Domain: DomainAbility, Policy: PolicyCommonOnly,
	},
	Rorschach: {
		// 同上：名义上属 typology，但检索只走 common。
		Code: Rorschach, Key: "rorschach", Name: "罗夏墨迹测验",
		Domain: DomainTypology, Policy: PolicyCommonOnly,
	},

	// ── 不接入增强（2）──────────────────────────────────────────
	Zootopia: {
		Code: Zootopia, Key: "zootopia", Name: "疯狂动物城性格测试",
		Domain: "", Policy: PolicyDisabled,
	},
	SBTI: {
		Code: SBTI, Key: "sbti", Name: "SBTI 傻大个性格测试",
		Domain: "", Policy: PolicyDisabled,
	},
}

// ── 查询辅助 ────────────────────────────────────────────────────

// LookupTest 取一个测评的元数据。第二个返回值为 false 表示未登记——
// 调用方应当据此**跳过增强**，而不是猜一个域出来。
func LookupTest(code int) (TestMeta, bool) {
	m, ok := TestRegistry[code]
	return m, ok
}

// RetrievalDomains 返回该测评应当检索的构念域，顺序即请求体 test.domains。
//
//   - full        → [本域, common]
//   - common_only → [common]
//   - disabled    → nil（调用方不得调用增强接口）
//   - 未登记      → nil
//
// 这是 Python 侧唯一需要的「分区信息」：它不再判断测评类型（ADR-0008）。
func RetrievalDomains(code int) []string {
	m, ok := TestRegistry[code]
	if !ok {
		return nil
	}
	switch m.Policy {
	case PolicyFull:
		return []string{string(m.Domain), string(DomainCommon)}
	case PolicyCommonOnly:
		return []string{string(DomainCommon)}
	default:
		return nil
	}
}

// CanEnhance 报告该测评是否接入检索增强。
func CanEnhance(code int) bool {
	m, ok := TestRegistry[code]
	if !ok {
		return false
	}
	return m.Policy != PolicyDisabled
}

// NeedsReferralForced 报告该测评是否必须强制 needs_referral。
// 只有 common_only 的测评会返回 true：它们的报告不得给出个性化解读。
func NeedsReferralForced(code int) bool {
	m, ok := TestRegistry[code]
	return ok && m.Policy == PolicyCommonOnly
}

// RegisteredCodes 返回全部已登记测评的编码，升序。
func RegisteredCodes() []int {
	codes := make([]int, 0, len(TestRegistry))
	for code := range TestRegistry {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	return codes
}

// DomainTestCounts 返回各构念域登记的测评数，用于与附录 B 对账。
// 只统计有 Domain 的条目（disabled 的两个不计入任何域）。
func DomainTestCounts() map[Domain]int {
	counts := make(map[Domain]int, len(AllDomains))
	for _, d := range AllDomains {
		counts[d] = 0
	}
	for _, m := range TestRegistry {
		if m.Domain != "" {
			counts[m.Domain]++
		}
	}
	return counts
}

// ── TODO（P0-5 前置）─────────────────────────────────────────────
//
// 本文件还缺**阈值表**：42 个测评的分档边界与量表上界目前散落在
// internal/logic/evaluation/*.go 里，共 103 处变量与数值字面量的比较
// （其中 89 处决定 level/title，14 处决定其他字段），另有 12 处写死的
// max_score 字面量与 5 处 len(data)*5 的动态上界。
//
// 为什么必须搬进来：rag-service/app/api/schemas.py 的 Measure 校验器
// **强制**每个 dimension 带 threshold、每个 factor 带 clinical_cutoff，
// 而 Go 侧目前 grep 不到任何 threshold/cutoff 字样——P0-5 的规范化层
// 无从填这个字段。搬完之后，这里就是阈值的**单一来源**，Python 契约、
// prompt 文案、语料口径三者才能对得上。
//
// 搬的时候要注意口径衔接：SCL-90 的因子切分点是 2.0（SCL90.go），
// 而 data/raw/symptom 的语料反复强调「总分和阳性项目数都是描述性指标，
// 不是诊断阈值」。两者不冲突，但报告必须把「系统分档」与「临床阈值」
// 明确区分开，否则会产出自相矛盾的解读。

// ── 测量项与阈值（P0-5）──────────────────────────────────────────
//
// 这是 rag-service/app/api/schemas.py 里 Measure 校验器要的那张表：每个
// dimension 必须带 threshold、每个 factor 必须带 clinical_cutoff，否则
// Python 直接 422。在此之前这些数字散落在 internal/logic/evaluation/*.go
// 里（103 处比较、12 处写死的 max），Go 侧根本拼不出一个合规的 measures[]。
//
// ⚠️ 两种阈值的口径**必须区分开**，否则报告会自相矛盾：
//   - Threshold（dimension）：量表内部的**分档点**，是"系统怎么给你归类"，
//     不是临床诊断线。MBTI 的 E/I 分界 12 就是这么来的（3 题 × 7 分 = 21，
//     取中点偏上）。
//   - ClinicalCutoff（factor）：症状量表的**临床参考线**（SCL-90 为 2.0），
//     语料里反复强调"因子分是描述性指标，不是诊断阈值"——报告里必须
//     把这两句话放在一起说，不能只报"你超过 2.0 了"。

// MeasureKind 对应 Python 契约里 measures[].kind。
type MeasureKind string

const (
	MeasureDimension MeasureKind = "dimension" // 人格/兴趣维度的两极
	MeasureFactor    MeasureKind = "factor"    // 症状量表的因子分
)

// MeasureSpec 是一个测评需要序列化进 measures[] 的一项。
type MeasureSpec struct {
	Key   string      // measures[].key
	Label string      // 中文名；同时用于检索 query 与 prompt
	Kind  MeasureKind // dimension / factor

	// Source 是 result_data 里的取值位置：
	//   - 普通字段：直接是 map 的键（如 MBTI 的 "extraversion"）
	//   - FromFactors 为真时：在 factors[] 里按 Label 匹配（如 SCL-90 的"抑郁"）
	Source      string
	FromFactors bool

	// Max 仅 dimension 用：该维度满分量（MBTI 的 E/S 是 21，T/J 是 14）。
	Max float64
	// Threshold 仅 dimension 用：分档点，score ≥ threshold 取 HighPole。
	Threshold float64
	// ClinicalCutoff 仅 factor 用：临床参考线（SCL-90 为 2.0）。
	ClinicalCutoff float64

	// HighPole / LowPole 仅 dimension 用：两极的代号。
	HighPole string
	LowPole  string
}

// TestMeasures 是各测评的测量项定义。
//
// ⚠️ 目前只登记了白名单里的 MBTI(13) 与 SCL-90(1)。DISC(28) 的评分代码
// **没有任何绝对阈值**——它只比四个维度谁最高来定 dominant_type。没有阈值就
// 拼不出合规的 Measure（Python 强制 dimension 带 threshold），所以 DISC 暂时
// 走「未登记」分支被跳过，而不是随便编一个分档点。要接 DISC 必须先回答
// 「D 多少分算高」，那是产品口径问题，不能由代码猜。
var TestMeasures = map[int][]MeasureSpec{
	MBTI: {
		{
			Key: "EI", Label: "外倾(E)/内倾(I)", Kind: MeasureDimension,
			Source: "extraversion", Max: 21, Threshold: 12,
			HighPole: "E", LowPole: "I",
		},
		{
			Key: "SN", Label: "感觉(S)/直觉(N)", Kind: MeasureDimension,
			Source: "sensing", Max: 21, Threshold: 12,
			HighPole: "S", LowPole: "N",
		},
		{
			Key: "TF", Label: "思考(T)/情感(F)", Kind: MeasureDimension,
			Source: "thinking", Max: 14, Threshold: 8,
			HighPole: "T", LowPole: "F",
		},
		{
			Key: "JP", Label: "判断(J)/感知(P)", Kind: MeasureDimension,
			Source: "judging", Max: 14, Threshold: 8,
			HighPole: "J", LowPole: "P",
		},
	},
	SCL90: {
		{Key: "somatization", Label: "躯体化", Kind: MeasureFactor, Source: "躯体化", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "obsessive", Label: "强迫症状", Kind: MeasureFactor, Source: "强迫症状", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "interpersonal", Label: "人际关系敏感", Kind: MeasureFactor, Source: "人际关系敏感", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "depression", Label: "抑郁", Kind: MeasureFactor, Source: "抑郁", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "anxiety", Label: "焦虑", Kind: MeasureFactor, Source: "焦虑", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "hostility", Label: "敌对", Kind: MeasureFactor, Source: "敌对", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "phobic", Label: "恐怖", Kind: MeasureFactor, Source: "恐怖", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "paranoid", Label: "偏执", Kind: MeasureFactor, Source: "偏执", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "psychoticism", Label: "精神病性", Kind: MeasureFactor, Source: "精神病性", FromFactors: true, ClinicalCutoff: 2.0},
		{Key: "other", Label: "其他", Kind: MeasureFactor, Source: "其他", FromFactors: true, ClinicalCutoff: 2.0},
	},
}

// MeasuresOf 返回该测评的测量项定义；未登记返回 nil。
//
// 调用方必须把 nil 当成「这个测评还不能接入增强」而不是「没有测量项也要发」：
// Python 侧 measures 有 min_length=1，空数组直接 422。
func MeasuresOf(code int) []MeasureSpec {
	return TestMeasures[code]
}
