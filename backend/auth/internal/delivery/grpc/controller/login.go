package controller

import (
	"context"
	"errors"
	"net"

	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/delivery/grpc"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/auth_service"
)

func (c *controller) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	const op = "Login"

	clientIPStr := grpc.GetClientIP(ctx)

	backoffSession, ok := c.backoff.GetIfExists(clientIPStr, backoffConfigLoginGroupID)

	if ok && !backoffSession.IsAllowed() {
		tryAfter := backoffSession.NextAllowedUntilSeconds()
		return nil, appErrors.Chainf(
			appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
			"%s.%s", c.pkg, op,
		)
	}

	clientIP := net.ParseIP(clientIPStr)

	authDTO, err := c.authFacade.Auth.LoginByEmail(
		ctx,
		req.Email,
		req.Password,
		clientIP,
	)
	if err != nil {
		if errors.Is(err, appErrors.ErrUnauthorized) {
			backoffSession = c.backoff.Get(clientIPStr, backoffConfigLoginGroupID)
			backoffSession.AddCounter()
			if backoffSession.Counter() > 1 {
				_ = backoffSession.AddBackoff()
			}
		}

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

	out := &desc.LoginResponse{
		SessionId: authDTO.Session.ID.String(),
		Tokens: &desc.Tokens{
			RefreshJwt: refreshJWT,
			AccessJwt:  accessJWT,
		},
		Account: &desc.Account{
			Id:             authDTO.AccountDTO.Account.ID.String(),
			Email:          authDTO.AccountDTO.Account.Email,
			ProfileName:    authDTO.AccountDTO.Account.ProfileName,
			ProfileSurname: authDTO.AccountDTO.Account.ProfileSurname,
		},
	}

	return out, nil
}
