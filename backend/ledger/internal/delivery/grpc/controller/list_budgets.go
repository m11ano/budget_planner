package controller

import (
	"context"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) ListBudgets(
	ctx context.Context,
	req *desc.ListBudgetsRequest,
) (*desc.ListBudgetsResponse, error) {
	const op = "ListBudgets"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	if req.Limit < 1 {
		req.Limit = 100
	}

	if req.Limit > 100 {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithHints("limit must be <= 100"), "%s.%s", c.pkg, op)
	}

	listOptions := &budgetUC.BudgetListOptions{
		Sort: []uctypes.SortOption[budgetUC.BudgetListOptionsSortField]{
			{
				Field:  budgetUC.BudgetListOptionsSortFieldPeriod,
				IsDesc: true,
			},
		},
		FilterAccountID: &authData.AccountID,
	}

	if req.FilterPeriodFrom != nil {
		occuredOn, err := civil.ParseDate(
			fmt.Sprintf(
				"%04d-%02d-%02d",
				req.FilterPeriodFrom.Year,
				req.FilterPeriodFrom.Month,
				1,
			),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid filter_period_from"), "%s.%s", c.pkg, op)
		}

		listOptions.FilterPeriodFrom = &occuredOn
	}

	if req.FilterPeriodTo != nil {
		occuredOn, err := civil.ParseDate(
			fmt.Sprintf(
				"%04d-%02d-%02d",
				req.FilterPeriodTo.Year,
				req.FilterPeriodTo.Month,
				1,
			),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid filter_period_to"), "%s.%s", c.pkg, op)
		}

		listOptions.FilterPeriodTo = &occuredOn
	}

	items, total, hitCache, err := c.budgetFacade.Budget.FindPagedList(
		ctx,
		listOptions,
		&uctypes.QueryGetListParams{
			Limit:  uint64(req.Limit),
			Offset: uint64(req.Offset),
		},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.ListBudgetsResponse{
		Items:    make([]*desc.Budget, 0, len(items)),
		Total:    int64(total),
		HitCache: hitCache,
	}

	for _, item := range items {
		out.Items = append(out.Items, BudgetToProto(item))
	}

	return out, nil
}
