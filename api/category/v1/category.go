package v1

import "github.com/gogf/gf/v2/frame/g"

type GetIndexReq struct {
	g.Meta `path:"/index" method:"get" tags:"Category" summary:"获取分类及测试列表"`
}

type TestTypeItem struct {
	Id          int64  `json:"id" dc:"测试类型ID"`
	Name        string `json:"name" dc:"测试类型名称"`
	Img         string `json:"img" dc:"测试图标"`
	Description string `json:"description" dc:"测试描述"`
	TestTime    int    `json:"testTime" dc:"预计测试时间"`
	UserCount   int    `json:"userCount" dc:"参与人数"` // 可选，模拟数据
}

type CategoryItem struct {
	Id        int64          `json:"id" dc:"分类ID"`
	Name      string         `json:"name" dc:"分类名称"`
	Sort      int            `json:"sort" dc:"排序优先级"`
	TestTypes []TestTypeItem `json:"testTypes" dc:"该分类下的测试"`
}

type GetIndexRes struct {
	List []CategoryItem `json:"list" dc:"分类列表"`
}
