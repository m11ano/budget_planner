package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) GetBudget(ctx context.Context, req *desc.GetBudgetRequest) (*desc.GetBudgetResponse, error) {
	const op = "GetBudget"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	budgetID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"), "%s.%s", c.pkg, op)
	}

	itemDTO, err := c.budgetFacade.Budget.FindOneByID(ctx, budgetID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.GetBudgetResponse{
		Item: BudgetToProto(itemDTO),
	}

	return out, nil
}
