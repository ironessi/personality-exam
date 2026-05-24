package pay

import (
	"context"

	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
)

type Price struct{}

func NewPrice() *Price {
	return &Price{}
}

type GetTestPriceInput struct {
	TestId int64
}

type GetTestPriceOutput struct {
	TestId        int64   `json:"testId"`
	TestName      string  `json:"testName"`
	OriginalPrice float64 `json:"originalPrice"`
	DiscountPrice float64 `json:"discountPrice"`
	DiscountRate  float64 `json:"discountRate"`
	PromotionText string  `json:"promotionText"`
}

func (p *Price) GetTestPrice(ctx context.Context, in GetTestPriceInput) (*GetTestPriceOutput, error) {
	var testType entity.TestTypes
	err := dao.TestTypes.Ctx(ctx).WherePri(in.TestId).Scan(&testType)
	if err != nil {
		return nil, err
	}
	if testType.Id == 0 {
		return nil, gerror.New("测试类型不存在")
	}

	originalPrice := testType.Price
	discountPrice := originalPrice
	discountRate := 1.0
	promotionText := ""

	return &GetTestPriceOutput{
		TestId:        testType.Id,
		TestName:      testType.Name,
		OriginalPrice: originalPrice,
		DiscountPrice: discountPrice,
		DiscountRate:  discountRate,
		PromotionText: promotionText,
	}, nil
}