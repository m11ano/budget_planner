package controller

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) GetTransaction(
	ctx context.Context,
	req *desc.GetTransactionRequest,
) (*desc.GetTransactionResponse, error) {
	const op = "GetTransaction"

	transactionID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"), "%s.%s", c.pkg, op)
	}

	itemDTO, err := c.budgetFacade.Transaction.FindOneByID(ctx, transactionID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.GetTransactionResponse{
		Item: TransactionToProto(itemDTO),
	}

	return out, nil
}
