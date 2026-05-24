package category

import (
	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/do"
	"PersonalityTest/internal/model/entity"
	"PersonalityTest/internal/service"
	"context"
)

type sCategory struct{}

func init() {
	service.RegisterCategory(New())
}

func New() *sCategory {
	return &sCategory{}
}

func (s *sCategory) GetTree(ctx context.Context) ([]*entity.TestCategories, map[int64][]*entity.TestTypes, error) {
	// 1. 获取所有启用状态的分类，按Sort排序
	var categories []*entity.TestCategories
	err := dao.TestCategories.Ctx(ctx).Where(do.TestCategories{Status: 1}).Order("sort asc").Scan(&categories)
	if err != nil {
		return nil, nil, err
	}

	// 2. 获取所有启用状态的测试
	var tests []*entity.TestTypes
	err = dao.TestTypes.Ctx(ctx).Where(do.TestTypes{Status: 1}).Scan(&tests)
	if err != nil {
		return nil, nil, err
	}

	// 3. 组装测试到Map: categoryId -> tests
	testMap := make(map[int64][]*entity.TestTypes)
	for _, t := range tests {
		// 这里假设 entity.TestTypes 已经添加了 CategoryId 字段
		// 如果 entity 还没有更新，这里会报错，所以前置步骤必须更新 entity
		cid := t.CategoryId
		if cid > 0 {
			testMap[cid] = append(testMap[cid], t)
		}
	}

	return categories, testMap, nil
}
