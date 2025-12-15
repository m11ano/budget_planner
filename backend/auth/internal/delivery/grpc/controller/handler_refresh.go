package controller

import (
	"context"
	"errors"
	"net"

	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/delivery/grpc"
	authUC "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/service"
)

func (c *controller) Refresh(ctx context.Context, req *desc.RefreshRequest) (*desc.RefreshResponse, error) {
	const op = "Refresh"

	clientIPStr := grpc.GetClientIP(ctx)
	clientIP := net.ParseIP(clientIPStr)

	claims, err := c.authFacade.Auth.ParseRefreshToken(req.RefreshJwt, true)
	if err != nil {
		if errors.Is(err, authUC.ErrInvalidToken) {
			return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", c.pkg, op)
		}
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	authDTO, err := c.authFacade.Auth.GenerateNewClaims(ctx, claims, clientIP)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	accessJWT, err := c.authFacade.Auth.IssueAccessJWT(authDTO.AccessClaims)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	refreshJWT, err := c.authFacade.Auth.IssueRefreshJWT(authDTO.RefreshClaims)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	out := &desc.RefreshResponse{
		Tokens: &desc.Tokens{
			RefreshJwt: refreshJWT,
			AccessJwt:  accessJWT,
		},
	}

	return out, nil
}
