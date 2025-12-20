package budget

import (
	"go.uber.org/fx"

	budgetRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/pg/budget"
	categoryRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/pg/category"
	transactionRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/pg/transaction"
	transactionRedisRepo "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/repository/redis/transaction"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase/budget"
	categoryUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase/category"
	transactionUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase/transaction"
)

var FxModule = fx.Module(
	"budget_module",

	// repositories
	fx.Provide(
		fx.Private,
		fx.Annotate(categoryRepo.NewRepository, fx.As(new(usecase.CategoryRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(transactionRepo.NewRepository, fx.As(new(usecase.TransactionRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(budgetRepo.NewRepository, fx.As(new(usecase.BudgetRepository))),
	),
	fx.Provide(
		fx.Private,
		fx.Annotate(transactionRedisRepo.NewRepository, fx.As(new(usecase.TransactionRedisRepository))),
	),

	// usecases
	fx.Provide(
		fx.Annotate(categoryUC.NewUsecaseImpl, fx.As(new(usecase.CategoryUsecase))),
	),
	fx.Provide(
		fx.Annotate(transactionUC.NewUsecaseImpl, fx.As(new(usecase.TransactionUsecase))),
	),
	fx.Provide(
		fx.Annotate(budgetUC.NewUsecaseImpl, fx.As(new(usecase.BudgetUsecase))),
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
