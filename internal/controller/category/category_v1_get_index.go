package category

import (
	v1 "PersonalityTest/api/category/v1"
	"PersonalityTest/internal/service"
	"context"
)

func (c *ControllerV1) GetIndex(ctx context.Context, req *v1.GetIndexReq) (res *v1.GetIndexRes, err error) {
	// 调用 Logic 获取数据
	categories, testMap, err := service.Category().GetTree(ctx)
	if err != nil {
		return nil, err
	}

	res = &v1.GetIndexRes{
		List: make([]v1.CategoryItem, 0),
	}

	for _, cat := range categories {
		item := v1.CategoryItem{
			Id:        cat.Id,
			Name:      cat.Name,
			Sort:      cat.Sort,
			TestTypes: make([]v1.TestTypeItem, 0),
		}

		if tests, ok := testMap[cat.Id]; ok {
			for _, t := range tests {
				item.TestTypes = append(item.TestTypes, v1.TestTypeItem{
					Id:          t.Id,
					Name:        t.Name,
					Img:         t.Img,
					Description: t.Description,
					TestTime:    t.TestTime,
				})
			}
		}

		res.List = append(res.List, item)
	}

	return res, nil
}
