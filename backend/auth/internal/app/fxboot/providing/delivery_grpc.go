package providing

import (
	"io"
	"log/slog"

	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	deliveryGRPC "github.com/m11ano/budget_planner/backend/auth/internal/delivery/grpc"
	controllerGRPC "github.com/m11ano/budget_planner/backend/auth/internal/delivery/grpc/controller"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var DeliveryGRPC = fx.Options(
	fx.Provide(func(logger *slog.Logger, cfg config.Config) *grpc.Server {
		if !cfg.BackendApp.GRPC.LogQueries {
			logger = slog.New(slog.NewTextHandler(io.Discard, nil))
		}
		return deliveryGRPC.New(logger, cfg)
	}),
	fx.Invoke(controllerGRPC.Register),
)
