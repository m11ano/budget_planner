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

func (c *controller) AddBudget(ctx context.Context, req *desc.AddBudgetRequest) (*desc.AddBudgetResponse, error) {
	const op = "AddBudget"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	amount, err := decimal.Parse(req.Amount)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid amount"), "%s.%s", c.pkg, op)
	}

	period, err := civil.ParseDate(
		fmt.Sprintf("%04d-%02d-%02d", req.Period.Year, req.Period.Month, 1),
	)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid period"), "%s.%s", c.pkg, op)
	}

	itemDTO, err := c.budgetFacade.Budget.CreateBudgetByDTO(
		ctx,
		budgetUC.CreateBudgetDataInput{
			AccountID:  authData.AccountID,
			Amount:     amount,
			Period:     period,
			CategoryID: uint64(req.CategoryId),
		},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.AddBudgetResponse{
		Item: &desc.Budget{
			Id:        itemDTO.Budget.ID.String(),
			AccountId: itemDTO.Budget.AccountID.String(),
			Amount:    itemDTO.Budget.Amount.String(),
			Period: &desc.DateMonth{
				Year:  int32(itemDTO.Budget.Period.Year),
				Month: int32(itemDTO.Budget.Period.Month),
			},
			CategoryId: int64(itemDTO.Budget.CategoryID),
			CreatedAt:  toProtoTimestamp(&itemDTO.Budget.CreatedAt),
			UpdatedAt:  toProtoTimestamp(&itemDTO.Budget.UpdatedAt),
		},
	}

	return out, nil
}
