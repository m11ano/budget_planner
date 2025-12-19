package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
)

func (c *controller) DeleteBudget(ctx context.Context, req *desc.DeleteBudgetRequest) (*desc.DeleteBudgetResponse, error) {
	const op = "DeleteBudget"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	budgetID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"), "%s.%s", c.pkg, op)
	}

	err = c.budgetFacade.Budget.DeleteBudgetByID(
		ctx,
		budgetID,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	return &desc.DeleteBudgetResponse{}, nil
}
