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

func (c *controller) CSVExportTransactions(
	ctx context.Context,
	req *desc.ListTransactionsRequest,
) (*desc.CSVExportTransactionsResponse, error) {
	const op = "CSVExportTransactions"

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

	listOptions := &budgetUC.TransactionListOptions{
		Sort: []uctypes.SortOption[budgetUC.TransactionListOptionsSortField]{
			{
				Field:  budgetUC.TransactionListOptionsSortFieldOccurredOn,
				IsDesc: true,
			},
			{
				Field:  budgetUC.TransactionListOptionsSortFieldCreatedAt,
				IsDesc: true,
			},
		},
		FilterAccountID: &authData.AccountID,
	}

	if req.FilterOccurredOnFrom != nil {
		occuredOn, err := civil.ParseDate(
			fmt.Sprintf(
				"%04d-%02d-%02d",
				req.FilterOccurredOnFrom.Year,
				req.FilterOccurredOnFrom.Month,
				req.FilterOccurredOnFrom.Day,
			),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid filter_occurred_on_from"), "%s.%s", c.pkg, op)
		}

		listOptions.FilterOccurredOnFrom = &occuredOn
	}

	if req.FilterOccurredOnTo != nil {
		occuredOn, err := civil.ParseDate(
			fmt.Sprintf(
				"%04d-%02d-%02d",
				req.FilterOccurredOnTo.Year,
				req.FilterOccurredOnTo.Month,
				req.FilterOccurredOnTo.Day,
			),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid filter_occurred_on_to"), "%s.%s", c.pkg, op)
		}

		listOptions.FilterOccurredOnTo = &occuredOn
	}

	data, total, err := c.budgetFacade.Transaction.FindPagedListAsCSV(
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

	out := &desc.CSVExportTransactionsResponse{
		Data:  data,
		Total: int64(total),
	}

	return out, nil
}
