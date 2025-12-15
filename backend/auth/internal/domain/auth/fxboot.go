package auth

import (
	"go.uber.org/fx"

	accountRepo "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/repository/pg/account"
	sessionRepo "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/repository/pg/session"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	accountUC "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase/account"
	authUC "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase/auth"
	sessionUC "github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase/session"
)

var FxModule = fx.Module(
	"tenant_module",

	// repositories
	fx.Provide(
		fx.Private,
		fx.Annotate(accountRepo.NewRepository, fx.As(new(usecase.AccountRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(sessionRepo.NewRepository, fx.As(new(usecase.SessionRepository))),
	),

	// usecases
	fx.Provide(
		fx.Annotate(accountUC.NewUsecaseImpl, fx.As(new(usecase.AccountUsecase))),
	),
	fx.Provide(
		fx.Annotate(authUC.NewUsecaseImpl, fx.As(new(usecase.AuthUsecase))),
	),
	fx.Provide(
		fx.Annotate(sessionUC.NewUsecaseImpl, fx.As(new(usecase.SessionUsecase))),
	),

	// facade
	fx.Provide(
		usecase.NewFacade,
	),

	// init
	fx.Provide(
		fx.Annotate(Init, fx.ResultTags(`group:"InvokeInit"`)),
	),
)
