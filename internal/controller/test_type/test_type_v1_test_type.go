package test_type

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "PersonalityTest/api/test_type/v1"
)

// AddTestType 添加测试类型
func (c *ControllerV1) AddTestType(ctx context.Context, req *v1.AddTestTypeReq) (res *v1.AddTestTypeRes, err error) {
	if req.Name == "" || req.Desc == "" || req.Status == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}

	id, err := c.types.AddTestType(ctx, req.Name, req.Desc, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.AddTestTypeRes{Id: id}, nil
}

// GetTest 获取测试题目 - 路径已改为 /api/tests/{id}
func (c *ControllerV1) GetTest(ctx context.Context, req *v1.GetTestReq) (res *v1.GetTestRes, err error) {
	if req.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}

	data, err := c.types.GetTest(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var questions []v1.Question
	for _, v := range data {
		var options []v1.Option

		// 直接解析JSON，不再需要复杂的修复逻辑
		if v.OptionContent != "" {
			err = json.Unmarshal([]byte(v.OptionContent), &options)
			if err != nil {
				// 如果解析失败，记录错误并返回空数组
				options = []v1.Option{}
			}
		}

		questions = append(questions, v1.Question{
			Id:       v.QuestionId,
			Content:  v.QuestionContent,
			Options:  options,
			Sequence: v.Sequence,
		})
	}

	// 确保questions不会是nil
	if questions == nil {
		questions = []v1.Question{}
	}

	return &v1.GetTestRes{Questions: questions}, nil
}

// SubmitTestAnswers 提交测试答案 - 注意：这个接口已迁移到 user_test 控制器
// 现在这个方法应该不再使用，或者修改为内部调用
//func (c *ControllerV1) SubmitTestAnswers(ctx context.Context, req *v1.SubmitTestAnswersReq) (res *v1.SubmitTestAnswersRes, err error) {
//	// 暂时使用模拟用户ID = 1（去掉认证后）
//	userID := int64(1)
//
//	var result interface{}
//	switch req.Type {
//	case consts.SCL90:
//		result, err = c.evaluation.SCL90(req.AnswerRecords)
//	case consts.GeneralLifeSatisfaction:
//		result, err = c.evaluation.GeneralLifeSatisfaction(req.AnswerRecords)
//	case consts.PerceivedStress:
//		result, err = c.evaluation.PerceivedStress(req.AnswerRecords)
//	case consts.PROMISAnxiety:
//		result, err = c.evaluation.PROMISAnxiety(req.AnswerRecords)
//	case consts.SocialDistress:
//		result, err = c.evaluation.SocialDistress(req.AnswerRecords)
//	case consts.UCLALoneliness:
//		result, err = c.evaluation.UCLALoneliness(req.AnswerRecords)
//	case consts.LifeMeaning:
//		result, err = c.evaluation.LifeMeaning(req.AnswerRecords)
//	case consts.Panas:
//		result, err = c.evaluation.Panas(req.AnswerRecords)
//	case consts.SelfEfficacy:
//		result, err = c.evaluation.SelfEfficacy(req.AnswerRecords)
//	case consts.PhysicalAggression:
//		result, err = c.evaluation.PhysicalAggression(req.AnswerRecords)
//	case consts.Hostility:
//		result, err = c.evaluation.Hostility(req.AnswerRecords)
//	default:
//		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	// 保存结果
//	err = c.result.SaveUserResult(ctx, userID, int64(req.Type), result)
//	if err != nil {
//		return nil, err
//	}
//
//	return &v1.SubmitTestAnswersRes{Result: result}, nil
//}

// GetTypeList 获取测试类型列表 - 路径已改为 /api/tests/types
func (c *ControllerV1) GetTypeList(ctx context.Context, req *v1.GetTypeListReq) (res *v1.GetTypeListRes, err error) {
	total, list, err := c.types.GetTypeList(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var data []v1.TestType
	for _, row := range list {
		var refs []v1.ReferenceItem
		if row.References != "" {
			if jsonErr := json.Unmarshal([]byte(row.References), &refs); jsonErr != nil {
				g.Log().Warningf(ctx, "反序列化 references 失败 [id=%d]: %v", row.Id, jsonErr)
				refs = []v1.ReferenceItem{}
			}
		}
		if refs == nil {
			refs = []v1.ReferenceItem{}
		}
		data = append(data, v1.TestType{
			Id:         row.Id,
			Name:       row.Name,
			Desc:       row.Description,
			References: refs,
		})
	}
	if data == nil {
		data = []v1.TestType{}
	}
	return &v1.GetTypeListRes{
		Total: int64(total),
		List:  data,
	}, nil
}
