package consts

const (
	JwtKey = "db03d23b03ec405793b38f10592a2f34"
)

// 测试类型常量
const (
	SCL90 = iota + 1
	GeneralLifeSatisfaction
	PerceivedStress
	PROMISAnxiety
	SocialDistress
	UCLALoneliness
	LifeMeaning        // 新增：生命意义感测试
	Panas              // 新增：积极情感消极情感量表
	SelfEfficacy       // 新增：自我效能感测试
	PhysicalAggression // 新增：身体攻击测试
	Hostility          // 新增：敌意测试
	// 职业测试
	RIASEC          // 霍兰德职业兴趣测试
	MBTI            // 16型人格测试
	Personalities   // Personalities人格测试
	PSTRI           // PSTRI压力知觉测试量表
	PhilPersonality // 菲尔人格十项测试题
	// 趣味测试
	STScale       // 社会地位指数测试ST量表
	AMS           // 成就动机测试(AMS)
	Zootopia      // 疯狂动物城测试性格
	StanfordBinet // 斯坦福-比奈智力量表
	Rorschach     // 罗夏墨迹测验
	// 新增心理测评问卷
	BDI                         // 贝克抑郁自评量表
	SDSAnxiety                  // 焦虑症自测量表(SDS)
	HCL32                       // 轻躁狂症状自评量表(HCL-32)
	MHLAQ                       // 心理健康素养量表
	TwentyNineDegreePersonality // 29度人格测试
	EMP                         // 同理心通情量表
	DISC                        // DISC性格测试
	EPQ                         // EPQ艾森克人格测试
	SocialStatusIndex           // 社会地位指数测试ST量表
	Machiavellianism            // 马基雅维利人格测试
	// 新增心理测评问卷
	SelfControlIndex    // 自我控力指数评估
	PassivePersonality  // 被动型人格测试
	TrueColors          // 色彩性格测试（True Colors）
	ObedienceTest       // 服从性测试
	MMPI                // 多相人格测试
	AnxietyScale        // 焦虑性测试 As量表
	EPPS                // 爱德华个人偏好测验
	HSP                 // 敏感度测试（HSP）
	NPI                 // 自恋人格量表（NPI）
	BemSexRoleInventory // 贝姆性别角色测试量表
	SBTI                // SBTI傻大个性格测试
)

// 响应状态码
const (
	CodeSuccess = 0
	CodeError   = -1
)

// 测试状态
const (
	TestStatusActive   = 1
	TestStatusInactive = 0
)
