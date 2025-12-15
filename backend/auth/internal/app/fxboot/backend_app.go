package fxboot

import (
	"context"
	"log/slog"
	"time"

	"github.com/m11ano/budget_planner/backend/auth/internal/app"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/fxboot/invoking"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/fxboot/providing"
	deliveryGRPC "github.com/m11ano/budget_planner/backend/auth/internal/delivery/grpc"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/db"
	"github.com/m11ano/budget_planner/backend/auth/pkg/backoff"
	"github.com/m11ano/budget_planner/backend/auth/pkg/pgclient"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"google.golang.org/grpc"
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
			ProvidingIDDBClients: fx.Provide(
				func(logger *slog.Logger, cfg config.Config, shutdown fx.Shutdowner) db.MasterClient {
					return providing.NewDBClients(
						cfg.Postgres.Master.DSN,
						cfg.BackendApp.Base.LogSQLQueries,
						logger,
						shutdown,
					)
				},
			),
			ProvidingIDBackoff:    fx.Provide(providing.NewBackoff),
			ProvidingGRPCServer:   providing.DeliveryGRPC,
			ProvidingIDAuthModule: auth.FxModule,
		},
		Invokes: []fx.Option{
			fx.Invoke(BackendAppInitInvoke),
		},
	}
}

type BAckendInvokeInput struct {
	fx.In

	LC             fx.Lifecycle
	Shutdowner     fx.Shutdowner
	Invokes        []invoking.InvokeInit `group:"InvokeInit"`
	Logger         *slog.Logger
	Cfg            config.Config
	DBMasterClient db.MasterClient
	BackoffCtrl    *backoff.Controller
	GRPCServer     *grpc.Server
}

func BackendAppInitInvoke(
	in BAckendInvokeInput,
) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		err := in.Shutdowner.Shutdown()
		if err != nil {
			in.Logger.Error("failed to shutdown", slog.Any("error", err))
		}
	}()

	in.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Тестирование соединения с мастером postgress
			err := pgclient.TestConnection(
				ctx,
				in.DBMasterClient,
				in.Logger,
				in.Cfg.Postgres.MaxAttempts,
				in.Cfg.Postgres.AttemptSleepSeconds,
			)
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to test master db connection", slog.Any("error", err))
				return err
			}

			in.Logger.InfoContext(
				ctx,
				"successfully connected to Postgress",
				slog.String("serverID", in.DBMasterClient.ServerID()),
			)

			// Миграции goose
			err = db.UpMigrations(in.Cfg.Postgres.Master.DSN, in.Cfg.Postgres.MigrationsPath, in.Logger)
			if err != nil {
				in.Logger.ErrorContext(ctx, "failed to run migrations", slog.Any("error", err))
				return err
			}

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
			if in.Cfg.BackendApp.GRPC.Port > 0 {
				in.Logger.InfoContext(ctx, "starting gRPC server", slog.Int("port", in.Cfg.BackendApp.GRPC.Port))
				go func() {
					err := deliveryGRPC.Start(in.GRPCServer, in.Cfg.BackendApp.GRPC.Port)
					if err != nil {
						in.Logger.ErrorContext(ctx, "failed to start gRPC server", slog.Any("error", err.Error()))
						cancel()
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

			// Останавливаем gRPC
			if in.Cfg.BackendApp.GRPC.Port > 0 {
				in.Logger.InfoContext(ctx, "stopping gRPC server")
				in.GRPCServer.GracefulStop()
			}

			// Закрываем postgress
			in.DBMasterClient.Close()
			in.Logger.InfoContext(ctx, "closing db clients")

			// Останавливаем backoff
			in.BackoffCtrl.Stop(ctx)

			return nil
		},
	})
}
