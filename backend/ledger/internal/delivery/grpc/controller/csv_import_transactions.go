package controller

import (
	"context"

	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) CSVImportTransactions(
	ctx context.Context,
	req *desc.CSVImportTransactionsRequest,
) (*desc.CSVImportTransactionsResponse, error) {
	const op = "CSVImportTransactions"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	err := c.budgetFacade.Transaction.ImportTransactionsFromCSV(
		ctx,
		req.Data,
		authData.AccountID,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	return &desc.CSVImportTransactionsResponse{}, nil
}
