package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) DeleteTransaction(
	ctx context.Context,
	req *desc.DeleteTransactionRequest,
) (*desc.DeleteTransactionResponse, error) {
	const op = "DeleteTransaction"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	transactionID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"), "%s.%s", c.pkg, op)
	}

	err = c.budgetFacade.Transaction.DeleteTransactionByID(
		ctx,
		transactionID,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	return &desc.DeleteTransactionResponse{}, nil
}
