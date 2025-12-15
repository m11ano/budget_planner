package controller

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/service"
)

func (c *controller) Whoiam(ctx context.Context, req *desc.WhoiamRequest) (*desc.WhoiamResponse, error) {
	const op = "Whoiam"

	sessionID, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", c.pkg, op)
	}

	// TODO: перенести в юзкейс метод!!
	session, err := c.authFacade.Session.FindOneByID(
		ctx,
		sessionID,
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	accountDTO, err := c.authFacade.Account.FindOneByID(
		ctx,
		session.AccountID,
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.WhoiamResponse{
		Account: &desc.Account{
			Id:             accountDTO.Account.ID.String(),
			Email:          accountDTO.Account.Email,
			ProfileName:    accountDTO.Account.ProfileName,
			ProfileSurname: accountDTO.Account.ProfileSurname,
		},
	}

	return out, nil
}
