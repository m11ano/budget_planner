package controller

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/auth_service"
)

func (c *controller) GetAccountByID(ctx context.Context, req *desc.GetAccountByIDRequest) (*desc.GetAccountByIDResponse, error) {
	const op = "GetAccountByID"

	id, err := uuid.Parse(req.AccountId)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", c.pkg, op)
	}

	accountDTO, err := c.authFacade.Account.FindOneByID(
		ctx,
		id,
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.GetAccountByIDResponse{
		Account: &desc.Account{
			Id:             accountDTO.Account.ID.String(),
			Email:          accountDTO.Account.Email,
			ProfileName:    accountDTO.Account.ProfileName,
			ProfileSurname: accountDTO.Account.ProfileSurname,
		},
	}

	return out, nil
}
