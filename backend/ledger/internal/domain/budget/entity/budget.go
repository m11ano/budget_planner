package entity

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
)

var ErrBudgetAmountInvalid = appErrors.ErrBadRequest.Extend("budget amount is invalid").
	WithTextCode("BUDGET_AMOUNT_INVALID").WithHints("budget amount is invalid")

type Budget struct {
	ID         uuid.UUID
	AccountID  uuid.UUID
	Period     civil.Date
	CategoryID uint64
	Amount     decimal.Decimal

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Budget) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Budget) SetAmount(value decimal.Decimal) error {
	if value.Cmp(decimal.Zero) == -1 {
		return ErrIncomeAmountInvalid
	}

	value = value.Trunc(2)

	item.Amount = value

	return nil
}

func (item *Budget) SetPeriod(value civil.Date) error {
	value.Day = 1
	item.Period = value

	return nil
}

func NewBudget(
	accountID uuid.UUID,
	amount decimal.Decimal,
	period civil.Date,
	categoryID uint64,
) (*Budget, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &Budget{
		ID:         uuid.New(),
		AccountID:  accountID,
		CategoryID: categoryID,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
	}

	err := item.SetAmount(amount)
	if err != nil {
		return nil, err
	}

	err = item.SetPeriod(period)
	if err != nil {
		return nil, err
	}

	return item, nil
}
