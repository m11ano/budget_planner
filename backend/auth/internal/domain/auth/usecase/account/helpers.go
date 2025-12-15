package account

import (
	"context"

	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
)

func (uc *UsecaseImpl) entitiesToDTO(
	_ context.Context,
	items []*entity.Account,
) ([]*usecase.AccountDTO, error) {
	out := make([]*usecase.AccountDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.AccountDTO{
			Account: item,
		}

		out = append(out, resItem)
	}

	return out, nil
}
