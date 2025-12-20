package transaction

import (
	"context"
	"fmt"
	"strings"

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
