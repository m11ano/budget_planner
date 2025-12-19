package controller

import (
	"context"

	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/auth_service"
)

func (c *controller) WhoIAm(ctx context.Context, req *desc.WhoIAmRequest) (*desc.WhoIAmResponse, error) {
	const op = "WhoIAm"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	accountDTO, err := c.authFacade.Account.FindOneByID(
		ctx,
		authData.AccountID,
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.WhoIAmResponse{
		Account: &desc.Account{
			Id:             accountDTO.Account.ID.String(),
			Email:          accountDTO.Account.Email,
			ProfileName:    accountDTO.Account.ProfileName,
			ProfileSurname: accountDTO.Account.ProfileSurname,
		},
	}

	return out, nil
}
