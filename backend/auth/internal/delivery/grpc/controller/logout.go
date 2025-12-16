package controller

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/auth_service"
)

func (c *controller) Logout(ctx context.Context, req *desc.LogoutRequest) (*desc.LogoutResponse, error) {
	const op = "Logout"

	id, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", c.pkg, op)
	}

	err = c.authFacade.Session.RevokeSessionByID(ctx, id)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	return &desc.LogoutResponse{}, nil
}
