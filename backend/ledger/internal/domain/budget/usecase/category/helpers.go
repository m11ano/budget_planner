package category

import (
	"context"

	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
)

func (uc *UsecaseImpl) entitiesToDTO(
	_ context.Context,
	items []*entity.Category,
) ([]*usecase.CategoryDTO, error) {
	out := make([]*usecase.CategoryDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.CategoryDTO{
			Category: item,
		}

		out = append(out, resItem)
	}

	return out, nil
}
