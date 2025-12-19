package controller

import (
	"context"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
	"github.com/samber/lo"
)

func (c *controller) PatchTransaction(
	ctx context.Context,
	req *desc.PatchTransactionRequest,
) (*desc.PatchTransactionResponse, error) {
	const op = "PatchTransaction"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	transactionID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"), "%s.%s", c.pkg, op)
	}

	patch := budgetUC.PatchTransactionDataInput{}

	if req.Amount != nil {
		amount, err := decimal.Parse(*req.Amount)
		if err != nil {
			return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid amount"), "%s.%s", c.pkg, op)
		}

		patch.Amount = &amount
	}

	if req.OccurredOn != nil {
		occuredOn, err := civil.ParseDate(
			fmt.Sprintf("%04d-%02d-%02d", req.OccurredOn.Year, req.OccurredOn.Month, req.OccurredOn.Day),
		)
		if err != nil {
			return nil, appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid occured_on"), "%s.%s", c.pkg, op)
		}

		patch.OccurredOn = &occuredOn
	}

	if req.CategoryId != nil {
		patch.CategoryID = lo.ToPtr(uint64(*req.CategoryId))
	}

	if req.Description != nil {
		patch.Description = lo.ToPtr(*req.Description)
	}

	err = c.budgetFacade.Transaction.PatchTransactionByDTO(
		ctx,
		transactionID,
		patch,
		true,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	itemDTO, err := c.budgetFacade.Transaction.FindOneByID(ctx, transactionID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.PatchTransactionResponse{
		Item: &desc.Transaction{
			Id:        itemDTO.Transaction.ID.String(),
			AccountId: itemDTO.Transaction.AccountID.String(),
			IsIncome:  itemDTO.Transaction.IsIncome,
			Amount:    itemDTO.Transaction.Amount.String(),
			OccurredOn: &desc.Date{
				Year:  int32(itemDTO.Transaction.OccurredOn.Year),
				Month: int32(itemDTO.Transaction.OccurredOn.Month),
				Day:   int32(itemDTO.Transaction.OccurredOn.Day),
			},
			CategoryId:  int64(itemDTO.Transaction.CategoryID),
			Description: itemDTO.Transaction.Description,
			CreatedAt:   toProtoTimestamp(&itemDTO.Transaction.CreatedAt),
			UpdatedAt:   toProtoTimestamp(&itemDTO.Transaction.UpdatedAt),
		},
	}

	return out, nil
}
