package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 添加测试类型 - 路径改为 /api/tests/types
type AddTestTypeReq struct {
	g.Meta `path:"/types" method:"post" tags:"TestType" summary:"添加测试类型"`
	Name   string `json:"name" v:"required#请输入测试类型名称" dc:"测试类型名称"`
	Desc   string `json:"desc" v:"required#请输入测试类型描述" dc:"测试类型描述"`
	Status int    `json:"status" v:"required#请输入测试类型状态" dc:"测试类型状态"`
}

type AddTestTypeRes struct {
	Id int64 `json:"id" dc:"测试类型ID"`
}

type Question struct {
	Id       int64    `json:"id" dc:"题目ID"`
	Sequence int64    `json:"sequence" dc:"题目序号"`
	Content  string   `json:"content" dc:"题目内容"`
	Options  []Option `json:"options" dc:"题目选项"`
}

type Option struct {
	Tag     string `json:"tag" dc:"选项标签"`
	Content string `json:"content" dc:"选项内容"`
	Score   int64  `json:"score" dc:"选项分数"`
}

// 获取测试题目 - 路径改为 /api/tests/{id}
type GetTestReq struct {
	g.Meta `path:"/{id}" method:"get" tags:"Test" summary:"获取测试题目"`
	Id     int64 `json:"id" v:"required" dc:"测试ID"`
}

type GetTestRes struct {
	Questions []Question `json:"questions" dc:"题目列表"`
}

// 获取测试类型列表 - 路径改为 /api/tests/types
type GetTypeListReq struct {
	g.Meta   `path:"/types" method:"get" tags:"TestType" summary:"获取测试类型列表"`
	Page     int `json:"page" v:"required#请输入页码" dc:"页码"`
	PageSize int `json:"pageSize" v:"required#请输入每页数量" dc:"每页数量"`
}

type ReferenceItem struct {
	Authors []string `json:"authors" dc:"作者列表"`
	Title   string   `json:"title"   dc:"文献标题"`
	Journal string   `json:"journal" dc:"期刊或出版社"`
	Year    int      `json:"year"    dc:"发表年份"`
	DOI     string   `json:"doi,omitempty" dc:"DOI标识符（可选）"`
}

type TestType struct {
	Id         int64           `json:"id"         dc:"测试类型ID"`
	Name       string          `json:"name"       dc:"测试类型名称"`
	Desc       string          `json:"desc"       dc:"测试类型描述"`
	References []ReferenceItem `json:"references" dc:"学术文献列表"`
}

type GetTypeListRes struct {
	Total int64      `json:"total" dc:"总记录数"`
	List  []TestType `json:"list" dc:"测试类型列表"`
}

type AnswerRecord struct {
	QuestionID int    `json:"questionId" v:"required#题目ID不能为空" dc:"题目序号"`
	Answer     string `json:"answer" v:"required#答案不能为空" dc:"用户作答内容"`
	Score      int64  `json:"score" v:"required#分数不能为空" dc:"用户作答得分"`
}

// SubmitTestAnswersReq 已迁移到 user_test 控制器，这里可以注释掉
//type SubmitTestAnswersReq struct {
//	g.Meta        `path:"/submit" method:"post" tags:"Test" summary:"提交心理测试作答"`
//	Type          int            `json:"type" query:"type" v:"required#测试类型不能为空" dc:"测试类型，如 mbti、depression 等"`
//	AnswerRecords []AnswerRecord `json:"answers" v:"" dc:"作答记录，每项包含题目ID和答案"`
//}

type SubmitTestAnswersRes struct {
	Result interface{} `json:"result" dc:"测试结果"`
}

type InitializeDataReq struct {
	g.Meta `path:"/initialize" method:"post" tags:"运维" summary:"初始化系统数据"`
	Force  bool `json:"force" dc:"是否强制重新初始化"`
}
type InitializeDataRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}

// 获取用户测试历史 - 这个应该使用新的 user_test 控制器
// 以下是旧的定义，建议使用新的 /api/user/tests 接口
type GetUserTestHistoryReq struct {
	g.Meta   `path:"/user/results" method:"get" tags:"Test" summary:"获取用户测试历史"`
	Page     int `json:"page" v:"min:1#页码最小为1" dc:"页码"`
	PageSize int `json:"pageSize" v:"min:1|max:100#每页数量范围1-100" dc:"每页数量"`
}

type UserTestResult struct {
	Id         int64       `json:"id" dc:"结果ID"`
	TestTypeId int64       `json:"testTypeId" dc:"测试类型ID"`
	TestName   string      `json:"testName" dc:"测试名称"`
	TestDate   *gtime.Time `json:"testDate" dc:"测试时间"`
	Level      string      `json:"level" dc:"结果等级"`
	Title      string      `json:"title" dc:"结果标题"`
	TotalScore int         `json:"totalScore" dc:"总分"`
}

type GetUserTestHistoryRes struct {
	Total int              `json:"total" dc:"总记录数"`
	List  []UserTestResult `json:"list" dc:"测试结果列表"`
}

// 获取特定测试结果详情 - 这个应该使用新的 user_test 控制器
// 以下是旧的定义，建议使用新的 /api/user/tests/{record_id} 接口
type GetTestResultDetailReq struct {
	g.Meta `path:"/user/results/{id}" method:"get" tags:"Test" summary:"获取测试结果详情"`
	Id     int64 `json:"id" v:"required" dc:"结果ID"`
}

type GetTestResultDetailRes struct {
	Result interface{} `json:"result" dc:"完整测试结果"`
}
