package transaction

import (
	"context"
	"fmt"
	"strings"

	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/entity"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) entitiesToDTO(
	ctx context.Context,
	items []*entity.Transaction,
) ([]*usecase.TransactionDTO, error) {
	out := make([]*usecase.TransactionDTO, 0, len(items))

	categoriesIDs := lo.Map(items, func(item *entity.Transaction, _ int) uint64 {
		return item.CategoryID
	})

	categories, err := uc.categoryRepo.FindList(ctx, &usecase.CategoryListOptions{
		FilterIDs: &categoriesIDs,
	}, nil)
	if err != nil {
		return nil, err
	}

	categoriesMap := lo.SliceToMap(categories, func(item *entity.Category) (uint64, *entity.Category) {
		return item.ID, item
	})

	for _, item := range items {
		resItem := &usecase.TransactionDTO{
			Transaction: item,
		}

		if category, ok := categoriesMap[item.CategoryID]; ok {
			resItem.Category = category
		}

		out = append(out, resItem)
	}

	return out, nil
}

func buildKeyForCountReportItems(queryFilter usecase.CountReportItemsQueryFilter) string {
	var strBuilder strings.Builder

	strBuilder.WriteString("Transaction::CountReportItems::")

	strBuilder.WriteString(queryFilter.AccountID.String())

	if queryFilter.DateFrom != nil {
		strBuilder.WriteString("::dateFrom:")
		strBuilder.WriteString(queryFilter.DateFrom.String())
	}

	if queryFilter.DateTo != nil {
		strBuilder.WriteString("::dateTo:")
		strBuilder.WriteString(queryFilter.DateTo.String())
	}

	if queryFilter.CategoryID != nil {
		strBuilder.WriteString("::categoryID:")
		strBuilder.WriteString(fmt.Sprintf("%d", *queryFilter.CategoryID))
	}

	return strBuilder.String()
}
