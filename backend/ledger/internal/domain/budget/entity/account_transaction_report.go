package entity

import (
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/samber/lo"
)

type AccountTransactionReportItem struct {
	Sum          *decimal.Decimal
	Period       civil.Date
	CategoryID   uint64
	BudgetID     *uuid.UUID
	BudgetAmount *decimal.Decimal
}

func (item *AccountTransactionReportItem) SpentBudget() (*decimal.Decimal, error) {
	if item.Sum == nil || item.BudgetAmount == nil || item.BudgetAmount.IsZero() {
		return nil, nil
	}

	if item.Sum.Cmp(decimal.Zero) >= 0 {
		return lo.ToPtr(decimal.Zero), nil
	}

	spentBudget, err := item.Sum.Neg().Quo(*item.BudgetAmount)
	if err != nil {
		return nil, err
	}

	spentBudgetPercents, err := spentBudget.Mul(decimal.Hundred)
	if err != nil {
		return nil, err
	}

	return lo.ToPtr(spentBudgetPercents.Trunc(2)), nil
}
