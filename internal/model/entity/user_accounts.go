// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/shopspring/decimal"
)

type UserAccounts struct {
	Id            int64           `json:"id" orm:"id"`
	UserId        int64          `json:"userId" orm:"user_id"`
	Balance       decimal.Decimal `json:"balance" orm:"balance"`
	FrozenAmount  decimal.Decimal `json:"frozenAmount" orm:"frozen_amount"`
	TotalRecharge decimal.Decimal `json:"totalRecharge" orm:"total_recharge"`
	TotalConsume  decimal.Decimal `json:"totalConsume" orm:"total_consume"`
	Version       int            `json:"version" orm:"version"`
}
