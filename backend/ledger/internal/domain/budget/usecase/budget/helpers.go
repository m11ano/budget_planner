package budget

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
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

func buildKeyForFindOneByID(id uuid.UUID) string {
	return fmt.Sprintf("Budget::FindOneByID::%s", id.String())
}

func buildKeyForFindList(
	listOptions *usecase.BudgetListOptions,
	queryParams *uctypes.QueryGetListParams,
) string {
	var strBuilder strings.Builder

	strBuilder.WriteString("Budget::FindList")

	if listOptions.FilterAccountID != nil {
		strBuilder.WriteString("::AccountID:")
		strBuilder.WriteString(listOptions.FilterAccountID.String())
	}

	if listOptions.FilterPeriod != nil {
		strBuilder.WriteString("::Period:")
		strBuilder.WriteString(listOptions.FilterPeriod.String())
	}

	if listOptions.FilterPeriodFrom != nil {
		strBuilder.WriteString("::PeriodFrom:")
		strBuilder.WriteString(listOptions.FilterPeriodFrom.String())
	}

	if listOptions.FilterPeriodTo != nil {
		strBuilder.WriteString("::PeriodTo:")
		strBuilder.WriteString(listOptions.FilterPeriodTo.String())
	}

	if listOptions.FilterCategoryID != nil {
		strBuilder.WriteString("::CategoryID:")
		strBuilder.WriteString(fmt.Sprintf("%d", *listOptions.FilterCategoryID))
	}

	if listOptions.Sort != nil {
		strBuilder.WriteString("::Sort:")
		for _, sort := range listOptions.Sort {
			strBuilder.WriteString(fmt.Sprintf("%s:%t", sort.Field, sort.IsDesc))
		}
	}

	if queryParams != nil {
		strBuilder.WriteString("::Offset:")
		strBuilder.WriteString(fmt.Sprintf("%d", queryParams.Offset))
		strBuilder.WriteString("::Limit:")
		strBuilder.WriteString(fmt.Sprintf("%d", queryParams.Limit))
	}

	return strBuilder.String()
}

func buildKeyForFindPagedList(
	listOptions *usecase.BudgetListOptions,
	queryParams *uctypes.QueryGetListParams,
) string {
	return fmt.Sprintf("%s::Paged", buildKeyForFindList(listOptions, queryParams))
}
