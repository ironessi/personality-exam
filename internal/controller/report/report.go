package report

import (
	"context"
	"encoding/json"

	"PersonalityTest/api/report/v1"
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/logic/result"
	"PersonalityTest/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type Controller struct {
	resultLogic *result.Result
}

func New() *Controller {
	return &Controller{
		resultLogic: result.New(),
	}
}

// GetDetail 获取报告详情
func (c *Controller) GetDetail(ctx context.Context, req *v1.GetReportReq) (res *v1.GetReportRes, err error) {
	userID := gconv.Int64(ctx.Value("uid"))
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized)
	}

	// 获取测试结果
	record, err := c.resultLogic.GetResultByID(ctx, req.RecordId, userID)
	if err != nil {
		return nil, gerror.New("查询报告失败")
	}

	if record == nil {
		return nil, gerror.New("报告不存在或无权访问")
	}

	// 获取测试类型信息
	var testType entity.TestTypes
	err = dao.TestTypes.Ctx(ctx).Where("id", record.TestTypeId).Scan(&testType)
	if err != nil {
		g.Log().Warningf(ctx, "获取测试类型失败: %v", err)
	}

	// 解析结果数据
	var resultData map[string]interface{}
	resultDataStr := gconv.String(record.ResultData)
	if err := json.Unmarshal([]byte(resultDataStr), &resultData); err != nil {
		return nil, gerror.New("解析分析结果失败")
	}

	// 构造返回
	res = &v1.GetReportRes{
		Id:         gconv.Int64(record.Id),
		TestId:     gconv.Int64(record.TestTypeId),
		TestName:   testType.Name,
		TotalScore: gconv.Int(resultData["total_score"]),
		Level:      gconv.String(resultData["level"]),
		Result:     resultData,
		CreatedAt:  gconv.GTime(record.TestDate),
	}

	return res, nil
}
