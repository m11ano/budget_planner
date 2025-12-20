package controller

import (
	"context"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/ledger/internal/common/uctypes"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
	"github.com/samber/lo"
)

func (c *controller) PatchBudget(ctx context.Context, req *desc.PatchBudgetRequest) (*desc.PatchBudgetResponse, error) {
	const op = "PatchBudget"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	budgetID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"), "%s.%s", c.pkg, op)
	}

	patch := budgetUC.PatchBudgetDataInput{}

	if req.Amount != nil {
		amount, err := decimal.Parse(*req.Amount)
		if err != nil {
			return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid amount"), "%s.%s", c.pkg, op)
		}

		patch.Amount = &amount
	}

	if req.Period != nil {
		period, err := civil.ParseDate(
			fmt.Sprintf("%04d-%02d-%02d", req.Period.Year, req.Period.Month, 1),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid period"), "%s.%s", c.pkg, op)
		}

		patch.Period = &period
	}

	if req.CategoryId != nil {
		patch.CategoryID = lo.ToPtr(uint64(*req.CategoryId))
	}

	err = c.budgetFacade.Budget.PatchBudgetByDTO(
		ctx,
		budgetID,
		patch,
		true,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	itemDTO, _, err := c.budgetFacade.Budget.FindOneByID(ctx, budgetID, &uctypes.QueryGetOneParams{
		SkipCache: true,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.PatchBudgetResponse{
		Item: BudgetToProto(itemDTO),
	}

	return out, nil
}
