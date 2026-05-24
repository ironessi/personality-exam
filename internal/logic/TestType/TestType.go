package TestType

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"context"
)

// TestTypeRow 扩展 entity.TestTypes，附加自动生成文件中未包含的 references 列。
type TestTypeRow struct {
	entity.TestTypes
	References string `json:"references" orm:"references"`
}

type TestType struct {
}

func New() *TestType {
	return &TestType{}
}

// AddTestType 添加测试类型(即一种测试)
func (t *TestType) AddTestType(ctx context.Context, name, desc string, status int) (id int64, err error) {
	result, err := dao.TestTypes.Ctx(ctx).Data(do.TestTypes{
		Name:        name,
		Description: desc,
		Status:      status,
	}).Insert()

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

type QuestionOption struct {
	QuestionId      int64  `json:"question_id" gconv:"question_id"`
	QuestionContent string `json:"question_content" gconv:"question_content"`
	Sequence        int64  `json:"sequence" gconv:"sequence"`
	OptionContent   string `json:"option_content" gconv:"option_content"`
}

func (t *TestType) GetTest(ctx context.Context, id int64) ([]QuestionOption, error) {
	// 首先获取所有问题
	var questions []struct {
		Id       int64  `json:"id"`
		Content  string `json:"content"`
		Sequence int64  `json:"sequence"`
	}

	err := dao.Questions.Ctx(ctx).
		Fields("id, content, sequence").
		Where("test_type_id", id).
		Order("sequence").
		Scan(&questions)
	if err != nil {
		return nil, err
	}

	var result []QuestionOption

	// 为每个问题获取选项
	for _, q := range questions {
		var optionRow struct {
			Content string `json:"content"`
		}

		// 直接从数据库获取选项JSON
		err := dao.Options.Ctx(ctx).
			Fields("content").
			Where("question_id", q.Id).
			Order("id").
			Scan(&optionRow)
		if err != nil {
			return nil, err
		}

		// ⚠️ 直接使用数据库中的JSON，不再重新构建
		result = append(result, QuestionOption{
			QuestionId:      q.Id,
			QuestionContent: q.Content,
			Sequence:        q.Sequence,
			OptionContent:   optionRow.Content, // 直接使用数据库的JSON
		})
	}

	return result, nil
}
func (t *TestType) GetTypeList(ctx context.Context, page, pageSize int) (int, []TestTypeRow, error) {
	var rows []TestTypeRow

	err := dao.TestTypes.Ctx(ctx).
		Fields("id, name, description, status, price, img, category_id, test_time, references").
		Where("status", 1).
		Page(page, pageSize).
		Scan(&rows)
	if err != nil {
		return 0, nil, err
	}

	count, err := dao.TestTypes.Ctx(ctx).Where("status", 1).Count()
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}
