package fxboot

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	authAdapter "github.com/m11ano/budget_planner/backend/gateway/internal/adapter/auth"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/config"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/fxboot/invoking"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/fxboot/providing"
	grpcClient "github.com/m11ano/budget_planner/backend/gateway/internal/client/grpc"
	backendHTTP "github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/backend"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func BackendAppGetOptionsMap(appID app.ID, cfg config.Config) OptionsMap {
	return OptionsMap{
		Providing: map[ProvidingID]fx.Option{
			ProvidingAppID: fx.Provide(func() app.ID {
				return appID
			}),
			ProvidingIDFXTimeouts: fx.Options(
				fx.StartTimeout(time.Second*time.Duration(cfg.BackendApp.Base.StartTimeoutSec)),
				fx.StopTimeout(time.Second*time.Duration(cfg.BackendApp.Base.StopTimeoutSec)),
			),
			ProvidingIDConfig: fx.Provide(func() config.Config {
				return cfg
			}),
			ProvidingIDLogger: fx.Provide(func(cfg config.Config) *slog.Logger {
				return providing.NewLogger(
					cfg.BackendApp.Name,
					cfg.BackendApp.Version,
					cfg.BackendApp.Base.UseLogger,
					cfg.BackendApp.Base.IsProd,
				)
			}),
			ProvidingIDFXLogger: fx.WithLogger(func(cfg config.Config) fxevent.Logger {
				return providing.NewFXLogger(cfg.BackendApp.Base.UseFxLogger)
			}),
			ProvidingHTTPFiberServer: fx.Provide(
				func(logger *slog.Logger, cfg config.Config) *fiber.App {
					httpConfig := backendHTTP.HTTPConfig{
						AppTitle:         cfg.BackendApp.Name,
						UnderProxy:       cfg.BackendApp.HTTP.UnderProxy,
						UseLogger:        cfg.BackendApp.Base.UseLogger && cfg.BackendApp.Base.LogHTTP,
						BodyLimit:        -1,
						CorsAllowOrigins: cfg.BackendApp.HTTP.CorsAllowOrigins,
						ServerIPs:        []string{cfg.Global.ServerIP},
					}

					return backendHTTP.NewHTTPFiber(httpConfig, logger)
				},
			),
			ProvidingIDDeliveryHTTP: backendHTTP.FxModule,
			ProvidingIDGRPCAuthClient: fx.Provide(
				func(logger *slog.Logger, cfg config.Config) authAdapter.Adapter {
					timeout, err := time.ParseDuration(cfg.GRPC.Auth.Timeout)
					if err != nil {
						panic(fmt.Errorf("failed to parse grpc auth timeout: %w", err))
					}

					adapter, err := authAdapter.NewAdapterImpl(cfg.GRPC.Auth.Addr, cfg.GRPC.Auth.RetriesCount, timeout, logger)
					if err != nil {
						panic(err)
					}

					return adapter
				},
			),
		},
		Invokes: []fx.Option{
			fx.Invoke(BackendAppInitInvoke),
		},
	}
}

type BAckendInvokeInput struct {
	fx.In

	LC              fx.Lifecycle
	Shutdowner      fx.Shutdowner
	Invokes         []invoking.InvokeInit `group:"InvokeInit"`
	Logger          *slog.Logger
	Cfg             config.Config
	HttpFiberServer *fiber.App
	AuthAdapter     authAdapter.Adapter
}

func BackendAppInitInvoke(
	in BAckendInvokeInput,
) {
	in.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Подключаем адаптеры
			err := grpcClient.ConnectToGRPCServer(ctx, in.AuthAdapter.CC())
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to connect to [auth] grpc server", slog.Any("error", err))
				return err
			}
			in.Logger.InfoContext(ctx, "connection established to [auth] grpc server")

			// Запускаем invoke функции до открытия
			for _, invokeItem := range in.Invokes {
				if invokeItem.StartBeforeOpen != nil {
					err := invokeItem.StartBeforeOpen(ctx)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn start before open", slog.Any("error", err))
						return err
					}
				}
			}

			// Запускаем http
			if in.Cfg.BackendApp.HTTP.Port > 0 {
				in.Logger.InfoContext(ctx, "starting http server", slog.Int("port", in.Cfg.BackendApp.HTTP.Port))
				go func() {
					if err := in.HttpFiberServer.Listen(fmt.Sprintf(":%d", in.Cfg.BackendApp.HTTP.Port)); err != nil {
						in.Logger.ErrorContext(ctx, "failed to start fiber", slog.Any("error", err))
						err := in.Shutdowner.Shutdown()
						if err != nil {
							in.Logger.ErrorContext(ctx, "failed to shutdown", slog.Any("error", err))
						}
					}
				}()
			}

			// Запускаем invoke функции после открытия
			for _, invokeItem := range in.Invokes {
				if invokeItem.StartAfterOpen != nil {
					err := invokeItem.StartAfterOpen(ctx)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn start after open", slog.Any("error", err))
						return err
					}
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			for _, invokeItem := range in.Invokes {
				if invokeItem.Stop != nil {
					err := invokeItem.Stop(ctx)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to execute invoke fn stop", slog.Any("error", err))
						return err
					}
				}
			}

			// Останавливаем http
			if in.Cfg.BackendApp.HTTP.Port > 0 {
				in.Logger.Info("stopping http fiber")
				err := in.HttpFiberServer.ShutdownWithTimeout(time.Duration(in.Cfg.BackendApp.HTTP.StopTimeoutSec) * time.Second)
				if err != nil {
					in.Logger.ErrorContext(ctx, "failed to stop fiber", slog.Any("error", err))
				}
			}

			// Отключаемся от auth
			err := in.AuthAdapter.CC().Close()
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to disconnect from [auth] grpc server", slog.Any("error", err))
				return err
			}

			return nil
		},
	})
}
