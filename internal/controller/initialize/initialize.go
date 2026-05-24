package initialize

import (
	"PersonalityTest/api/test_type/v1"

	"PersonalityTest/internal/logic/initialData"

	"PersonalityTest/internal/logic/TestType"
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type Controller struct {
	testType *TestType.TestType
}

func New() *Controller {
	return &Controller{
		testType: TestType.New(),
	}
}

func (c *Controller) InitializeData(ctx context.Context, req *v1.InitializeDataReq) (res *v1.InitializeDataRes, err error) {
	res = &v1.InitializeDataRes{}

	// 检查是否已存在测试类型数据
	count, err := c.checkExistingData(ctx)
	if err != nil {
		return nil, gerror.Newf("检查数据失败: %v", err)
	}

	// 如果已有数据且不强制初始化
	if count > 0 && !req.Force {
		res.Success = false
		res.Message = "数据库中已存在测试数据，如需重新初始化请设置 force=true"
		return res, nil
	}

	// 执行初始化
	g.Log().Info(ctx, "开始初始化测试数据...")
	initialData.ImportData()

	// 重新检查数据数量
	newCount, err := c.checkExistingData(ctx)
	if err != nil {
		return nil, gerror.Newf("检查初始化结果失败: %v", err)
	}

	res.Success = true
	res.Message = "数据初始化完成"
	res.Count = newCount

	g.Log().Info(ctx, "数据初始化完成，共初始化", newCount, "种测试")
	return res, nil
}

// checkExistingData 检查数据库中是否已存在测试数据
func (c *Controller) checkExistingData(ctx context.Context) (int, error) {
	// 检查 test_types 表中是否有数据
	count, err := g.Model("test_types").Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
