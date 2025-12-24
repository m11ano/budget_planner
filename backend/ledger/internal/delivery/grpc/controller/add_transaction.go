package controller

import (
	"context"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) AddTransaction(ctx context.Context, req *desc.AddTransactionRequest) (*desc.AddTransactionResponse, error) {
	const op = "AddTransaction"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	amount, err := decimal.Parse(req.Amount)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid amount"), "%s.%s", c.pkg, op)
	}

	occuredOn, err := civil.ParseDate(
		fmt.Sprintf("%04d-%02d-%02d", req.OccurredOn.Year, req.OccurredOn.Month, req.OccurredOn.Day),
	)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid occured_on"), "%s.%s", c.pkg, op)
	}

	itemDTO, err := c.budgetFacade.Transaction.CreateTransactionByDTO(
		ctx,
		budgetUC.CreateTransactionDataInput{
			AccountID:   authData.AccountID,
			IsIncome:    req.IsIncome,
			Amount:      amount,
			OccurredOn:  occuredOn,
			CategoryID:  uint64(req.CategoryId),
			Description: req.Description,
		},
		false,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.AddTransactionResponse{
		Item: TransactionToProto(itemDTO),
	}

	return out, nil
}
