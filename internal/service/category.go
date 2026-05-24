// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can use "gf gen service" to generate this file.
// ================================================================================

package service

import (
	"PersonalityTest/internal/model/entity"
	"context"
)

type (
	ICategory interface {
		GetTree(ctx context.Context) ([]*entity.TestCategories, map[int64][]*entity.TestTypes, error)
	}
)

var (
	localCategory ICategory
)

func Category() ICategory {
	if localCategory == nil {
		panic("implement not found for interface ICategory, forgot register?")
	}
	return localCategory
}

func RegisterCategory(i ICategory) {
	localCategory = i
}
