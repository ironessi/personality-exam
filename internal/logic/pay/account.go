package pay

import (
	"context"

	"PersonalityTest/internal/dao"
	"PersonalityTest/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shopspring/decimal"
)

type Account struct{}

func NewAccount() *Account {
	return &Account{}
}

type GetAccountOutput struct {
	UserId       int64           `json:"userId"`
	Balance      decimal.Decimal `json:"balance"`
	FrozenAmount decimal.Decimal `json:"frozenAmount"`
	TotalRecharge decimal.Decimal `json:"totalRecharge"`
	TotalConsume decimal.Decimal `json:"totalConsume"`
}

func (a *Account) GetAccount(ctx context.Context, userId int64) (*GetAccountOutput, error) {
	accountId, err := dao.UserAccounts.GetOrCreateAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	var account entity.UserAccounts
	err = dao.UserAccounts.Ctx(ctx).WherePri(accountId).Scan(&account)
	if err != nil {
		return nil, err
	}
	if account.Id == 0 {
		return nil, gerror.New("用户账户不存在")
	}

	return &GetAccountOutput{
		UserId:       account.UserId,
		Balance:      account.Balance,
		FrozenAmount: account.FrozenAmount,
		TotalRecharge: account.TotalRecharge,
		TotalConsume: account.TotalConsume,
	}, nil
}

type DeductBalanceInput struct {
	UserId   int64
	Amount   decimal.Decimal
	Source   string
	SourceId string
}

type DeductBalanceOutput struct {
	Balance    decimal.Decimal `json:"balance"`
	Deducted  decimal.Decimal `json:"deducted"`
}

func (a *Account) DeductBalance(ctx context.Context, in DeductBalanceInput) (*DeductBalanceOutput, error) {
	if in.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, gerror.New("扣减金额必须大于0")
	}

	var account entity.UserAccounts
	err := dao.UserAccounts.Ctx(ctx).
		Where("user_id", in.UserId).
		Scan(&account)
	if err != nil {
		return nil, err
	}
	if account.Id == 0 {
		return nil, gerror.New("用户账户不存在")
	}

	if account.Balance.LessThan(in.Amount) {
		return nil, gerror.New("余额不足")
	}

	newBalance := account.Balance.Sub(in.Amount)
	newVersion := account.Version + 1

	affected, err := dao.UserAccounts.Ctx(ctx).
		Where("id", account.Id).
		Where("version", account.Version).
		Data(g.Map{
			"balance":        newBalance,
			"total_consume":  account.TotalConsume.Add(in.Amount),
			"version":        newVersion,
		}).Update()
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := affected.RowsAffected()
	if rowsAffected == 0 {
		return nil, gerror.New("并发扣减失败，请重试")
	}

	return &DeductBalanceOutput{
		Balance:   newBalance,
		Deducted: in.Amount,
	}, nil
}

type AddBalanceInput struct {
	UserId  int64
	Amount  decimal.Decimal
	Source  string
}

type AddBalanceOutput struct {
	Balance   decimal.Decimal `json:"balance"`
	Added    decimal.Decimal `json:"added"`
}

func (a *Account) AddBalance(ctx context.Context, in AddBalanceInput) (*AddBalanceOutput, error) {
	if in.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, gerror.New("增加金额必须大于0")
	}

	accountId, err := dao.UserAccounts.GetOrCreateAccount(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	err = dao.UserAccounts.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var account entity.UserAccounts
		err := dao.UserAccounts.Ctx(ctx).
			WherePri(accountId).
			LockUpdate().
			Scan(&account)
		if err != nil {
			return err
		}

		newBalance := account.Balance.Add(in.Amount)
		newVersion := account.Version + 1

		_, err = dao.UserAccounts.Ctx(ctx).
			WherePri(accountId).
			Data(g.Map{
				"balance":         newBalance,
				"total_recharge":  account.TotalRecharge.Add(in.Amount),
				"version":         newVersion,
			}).Update()
		return err
	})
	if err != nil {
		return nil, err
	}

	var finalAccount entity.UserAccounts
	dao.UserAccounts.Ctx(ctx).WherePri(accountId).Scan(&finalAccount)

	return &AddBalanceOutput{
		Balance: finalAccount.Balance,
		Added:  in.Amount,
	}, nil
}