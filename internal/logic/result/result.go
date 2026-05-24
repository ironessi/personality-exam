package result

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"context"
	"encoding/json"
)

type Result struct{}

func New() *Result {
	return &Result{}
}

// SaveUserResult 保存用户测试结果到 user_test_results 表
func (r *Result) SaveUserResult(ctx context.Context, userID, testTypeID int64, result interface{}) error {
	resultByte, err := json.Marshal(result)
	if err != nil {
		return err
	}

	_, err = dao.UserTestResults.Ctx(ctx).Data(do.UserTestResults{
		UserId:     userID,
		TestTypeId: testTypeID,
		ResultData: string(resultByte),
	}).Insert()
	return err
}

// GetUserResults 获取用户的历史测试结果
func (r *Result) GetUserResults(ctx context.Context, userID int64, page, pageSize int) ([]do.UserTestResults, int, error) {
	var results []do.UserTestResults

	// 查询数据
	model := dao.UserTestResults.Ctx(ctx).
		Where("user_id", userID).
		OrderDesc("test_date")

	// 获取总数
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = model.Page(page, pageSize).Scan(&results)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetUserResultByTestType 获取用户特定测试类型的最新结果
func (r *Result) GetUserResultByTestType(ctx context.Context, userID, testTypeID int64) (*do.UserTestResults, error) {
	var result do.UserTestResults

	err := dao.UserTestResults.Ctx(ctx).
		Where("user_id", userID).
		Where("test_type_id", testTypeID).
		OrderDesc("test_date").
		Limit(1).
		Scan(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetResultByID 通过 ID 获取特定的测试结果
func (r *Result) GetResultByID(ctx context.Context, resultID, userID int64) (*do.UserTestResults, error) {
	var result do.UserTestResults

	err := dao.UserTestResults.Ctx(ctx).
		Where("id", resultID).
		Where("user_id", userID).
		Scan(&result)

	if err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, nil
	}

	return &result, nil
}

// DeleteUserResult 删除用户的测试结果（软删除或物理删除）
func (r *Result) DeleteUserResult(ctx context.Context, resultID, userID int64) error {
	_, err := dao.UserTestResults.Ctx(ctx).
		Where("id", resultID).
		Where("user_id", userID).
		Delete()
	return err
}
