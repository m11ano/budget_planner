package budget

import (
	"context"

	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
)

func (uc *UsecaseImpl) entitiesToDTO(
	_ context.Context,
	items []*entity.Budget,
) ([]*usecase.BudgetDTO, error) {
	out := make([]*usecase.BudgetDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.BudgetDTO{
			Budget: item,
		}

		out = append(out, resItem)
	}

	return out, nil
}
