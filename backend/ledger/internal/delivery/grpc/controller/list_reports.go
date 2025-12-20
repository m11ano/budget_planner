package controller

import (
	"context"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) ListReports(
	ctx context.Context,
	req *desc.ListReportsRequest,
) (*desc.ListReportsResponse, error) {
	const op = "ListReports"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	queryFilter := budgetUC.CountReportItemsQueryFilter{
		AccountID: authData.AccountID,
	}

	if req.DateFrom != nil {
		date, err := civil.ParseDate(
			fmt.Sprintf(
				"%04d-%02d-%02d",
				req.DateFrom.Year,
				req.DateFrom.Month,
				req.DateFrom.Day,
			),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid date_from"), "%s.%s", c.pkg, op)
		}

		queryFilter.DateFrom = &date
	}

	if req.DateTo != nil {
		date, err := civil.ParseDate(
			fmt.Sprintf(
				"%04d-%02d-%02d",
				req.DateTo.Year,
				req.DateTo.Month,
				req.DateTo.Day,
			),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid date_to"), "%s.%s", c.pkg, op)
		}

		queryFilter.DateTo = &date
	}

	items, err := c.budgetFacade.Transaction.CountReportItems(
		ctx,
		queryFilter,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.ListReportsResponse{
		Reports: make([]*desc.PeriodReport, 0, len(items)),
	}

	for _, item := range items {
		out.Reports = append(out.Reports, ReportToProto(item))
	}

	return out, nil
}
