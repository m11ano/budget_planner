package transaction

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	usecasemocks "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase/mocks"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/pgclient"
)

type dependencies struct {
	mc *minimock.Controller

	uc *UsecaseImpl

	logger *slog.Logger
	cfg    config.Config

	dbMasterClient       any
	transactionRepo      *usecasemocks.TransactionRepositoryMock
	transactionCacheRepo *usecasemocks.TransactionCacheRepositoryMock
	transactionCSVRepo   *usecasemocks.TransactionCSVRepositoryMock
	categoryRepo         *usecasemocks.CategoryRepositoryMock
	budgetRepo           *usecasemocks.BudgetRepositoryMock
}

func newDependencies(t *testing.T) *dependencies {
	t.Helper()

	mc := minimock.NewController(t)

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	transactionRepo := usecasemocks.NewTransactionRepositoryMock(mc)
	transactionCacheRepo := usecasemocks.NewTransactionCacheRepositoryMock(mc)
	transactionCSVRepo := usecasemocks.NewTransactionCSVRepositoryMock(mc)
	categoryRepo := usecasemocks.NewCategoryRepositoryMock(mc)
	budgetRepo := usecasemocks.NewBudgetRepositoryMock(mc)

	dbMasterClient := pgclient.NewMock()

	uc := NewUsecaseImpl(
		logger,
		config.Config{},
		dbMasterClient,
		transactionRepo,
		transactionCacheRepo,
		transactionCSVRepo,
		categoryRepo,
		budgetRepo,
	)

	return &dependencies{
		mc:                   mc,
		uc:                   uc,
		logger:               logger,
		cfg:                  config.Config{},
		dbMasterClient:       dbMasterClient,
		transactionRepo:      transactionRepo,
		transactionCacheRepo: transactionCacheRepo,
		transactionCSVRepo:   transactionCSVRepo,
		categoryRepo:         categoryRepo,
		budgetRepo:           budgetRepo,
	}
}

func finishDependencies(s *dependencies) {
	s.mc.Finish()
}

func testCtx() context.Context {
	return context.Background()
}
