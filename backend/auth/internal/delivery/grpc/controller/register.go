package controller

import (
	"context"
	"net"

	appErrors "github.com/m11ano/budget_planner/backend/auth/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/auth/internal/delivery/grpc"
	authUC "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/auth_service"
)

func (c *controller) Register(ctx context.Context, req *desc.RegisterRequest) (*desc.RegisterResponse, error) {
	const op = "Register"

	clientIPStr := grpc.GetClientIP(ctx)

	backoffSession, ok := c.backoff.GetIfExists(clientIPStr, backoffConfigRegisterGroupID)

	if ok && !backoffSession.IsAllowed() {
		tryAfter := backoffSession.NextAllowedUntilSeconds()

		return nil, appErrors.Chainf(
			appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
			"%s.%s", c.pkg, op,
		)
	}

	clientIP := net.ParseIP(clientIPStr)

	_, err := c.authFacade.Account.CreateAccountByDTO(
		ctx,
		authUC.CreateAccountDataInput{
			Email:          req.Email,
			Password:       req.Password,
			IsConfirmed:    true,
			ProfileName:    req.ProfileName,
			ProfileSurname: req.ProfileSurname,
		},
		&clientIP,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", c.pkg, op)
	}

	backoffSession = c.backoff.Get(clientIPStr, backoffConfigLoginGroupID)
	backoffSession.AddCounter()
	_ = backoffSession.AddBackoff()

	return &desc.RegisterResponse{}, nil
}
