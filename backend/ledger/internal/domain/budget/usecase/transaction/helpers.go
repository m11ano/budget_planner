package transaction

import (
	"context"

	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
)

func (uc *UsecaseImpl) entitiesToDTO(
	_ context.Context,
	items []*entity.Transaction,
) ([]*usecase.TransactionDTO, error) {
	out := make([]*usecase.TransactionDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.TransactionDTO{
			Transaction: item,
		}

		out = append(out, resItem)
	}

	return out, nil
}
