package controller

import (
	"time"

	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	authUC "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/pkg/backoff"
	desc "github.com/m11ano/budget_planner/backend/auth/pkg/proto_pb/auth_service"
	"google.golang.org/grpc"
)

type controller struct {
	desc.UnimplementedAuthServer
	pkg        string
	cfg        config.Config
	backoff    *backoff.Controller
	authFacade *authUC.Facade
}

const (
	backoffConfigLoginGroupID    = "grpc.login"
	backoffConfigRegisterGroupID = "grpc.register"
)

func Register(
	gRPCServer *grpc.Server,
	cfg config.Config,
	backoffCtrl *backoff.Controller,
	authFacade *authUC.Facade,
) {
	ctrl := &controller{
		pkg:        "grpc.Controller",
		cfg:        cfg,
		backoff:    backoffCtrl,
		authFacade: authFacade,
	}

	desc.RegisterAuthServer(gRPCServer, ctrl)

	ctrl.backoff.SetConfigForGroup(
		backoffConfigLoginGroupID,
		backoff.WithTtl(time.Minute*10),
		backoff.WithInitialInterval(time.Second*5),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*1),
	)

	ctrl.backoff.SetConfigForGroup(
		backoffConfigRegisterGroupID,
		backoff.WithTtl(time.Minute*10),
		backoff.WithInitialInterval(time.Second*5),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*1),
	)
}
