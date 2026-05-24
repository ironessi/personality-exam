package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type GetReportReq struct {
	g.Meta   `path:"/{RecordId}" method:"get" summary:"获取测试分析报告" tags:"报告"`
	RecordId int64 `json:"recordId" v:"required|min:1#记录ID必填|记录ID无效"`
}

type GetReportRes struct {
	Id         int64                  `json:"id" dc:"记录ID"`
	TestId     int64                  `json:"test_id" dc:"测试类型ID"`
	TestName   string                 `json:"test_name" dc:"测试名称"`
	TotalScore int                    `json:"total_score" dc:"总分"`
	Level      string                 `json:"level" dc:"等级"`
	Result     map[string]interface{} `json:"result" dc:"分析结果数据"`
	CreatedAt  *gtime.Time            `json:"created_at" dc:"测试时间"`
}
