package model

import (
	"PersonalityTest/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

// 获取测试题目请求
type GetTestQuestionsReq struct {
	TestTypeID int64 `v:"required|min:1#测试类型ID不能为空|测试类型ID无效" json:"test_type_id"`
}

type GetTestQuestionsRes struct {
	Questions []*entity.Questions `json:"questions"`
}

// 提交测试答案请求
type SubmitTestAnswersReq struct {
	TestTypeID int64               `v:"required|min:1#测试类型ID不能为空|测试类型ID无效" json:"test_type_id"`
	Answers    []entity.UserAnswers `json:"answers"`
}

type SubmitTestAnswersRes struct {
	Result   *entity.TestResults `json:"result"`
	ResultID int64                  `json:"result_id"`
}

// 测试类型相关请求
type CreateTestTypeReq struct {
	Name   string `v:"required#测试类型名称不能为空" json:"name"`
	Desc   string `json:"desc"`
	Status int    `json:"status"`
}

type CreateTestTypeRes struct {
	Message string `json:"message"`
}

type GetTestTypesReq struct {
	Page     int `d:"1"  v:"min:1#页码最小为1" json:"page"`
	PageSize int `d:"20" v:"max:100#每页数量最多100" json:"pageSize"`
}

type GetTestTypesRes struct {
	Total int                `json:"total"`
	List  []*entity.TestTypes `json:"list"`
}

// 初始化数据请求
type InitializeDataReq struct {
	g.Meta `path:"/initialize" method:"post" tags:"数据管理" summary:"初始化系统数据"`
	Force  bool `json:"force" dc:"是否强制重新初始化"`
}

type InitializeDataRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"`
}
