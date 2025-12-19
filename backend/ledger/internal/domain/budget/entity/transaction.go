package entity

import (
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
)

var (
	ErrIncomeAmountInvalid = appErrors.ErrBadRequest.Extend("income amount is invalid").
				WithTextCode("INCOME_AMOUNT_INVALID").WithHints("income amount is invalid")

	ErrOutcomeAmountInvalid = appErrors.ErrBadRequest.Extend("outcome amount is invalid").
				WithTextCode("OUTCOME_AMOUNT_INVALID").WithHints("outcome amount is invalid")
)

type Transaction struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	IsIncome    bool
	Amount      decimal.Decimal
	OccurredOn  civil.Date
	CategoryID  uint64
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (item *Transaction) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func (item *Transaction) SetDescription(value string) error {
	value = strings.TrimSpace(value)

	item.Description = value

	return nil
}

func (item *Transaction) SetAmount(value decimal.Decimal) error {
	if item.IsIncome && value.Cmp(decimal.Zero) == -1 {
		return ErrIncomeAmountInvalid
	} else if !item.IsIncome && value.Cmp(decimal.Zero) == 1 {
		return ErrOutcomeAmountInvalid
	}

	value = value.Trunc(2)

	item.Amount = value

	return nil
}

func (item *Transaction) SetOccuredOn(value civil.Date) error {
	item.OccurredOn = value

	return nil
}

func NewTransaction(
	accountID uuid.UUID,
	isIncome bool,
	amount decimal.Decimal,
	occurredOn civil.Date,
	categoryID uint64,
) (*Transaction, error) {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &Transaction{
		ID:         uuid.New(),
		AccountID:  accountID,
		IsIncome:   isIncome,
		CategoryID: categoryID,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
	}

	err := item.SetAmount(amount)
	if err != nil {
		return nil, err
	}

	err = item.SetOccuredOn(occurredOn)
	if err != nil {
		return nil, err
	}

	return item, nil
}
