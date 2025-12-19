package controller

import (
	"context"

	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/auth_service"
)

func (c *controller) Logout(ctx context.Context, req *desc.LogoutRequest) (*desc.LogoutResponse, error) {
	const op = "Logout"

	authData := auth.GetAuthData(ctx)
	if authData == nil {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
	}

	err := c.authFacade.Session.RevokeSessionByID(ctx, authData.SessionID)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	return &desc.LogoutResponse{}, nil
}
