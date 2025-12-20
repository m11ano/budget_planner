package budget

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

	dbMasterClient  any
	budgetRepo      *usecasemocks.BudgetRepositoryMock
	budgetCacheRepo *usecasemocks.BudgetCacheRepositoryMock
	categoryRepo    *usecasemocks.CategoryRepositoryMock
}

func newDependencies(t *testing.T) *dependencies {
	t.Helper()

	mc := minimock.NewController(t)

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	budgetRepo := usecasemocks.NewBudgetRepositoryMock(mc)
	budgetCacheRepo := usecasemocks.NewBudgetCacheRepositoryMock(mc)
	categoryRepo := usecasemocks.NewCategoryRepositoryMock(mc)

	dbMasterClient := pgclient.NewMock()

	uc := NewUsecaseImpl(
		logger,
		config.Config{},
		dbMasterClient,
		budgetRepo,
		budgetCacheRepo,
		categoryRepo,
	)

	return &dependencies{
		mc:              mc,
		uc:              uc,
		logger:          logger,
		cfg:             config.Config{},
		dbMasterClient:  dbMasterClient,
		budgetRepo:      budgetRepo,
		budgetCacheRepo: budgetCacheRepo,
		categoryRepo:    categoryRepo,
	}
}

func finishDependencies(s *dependencies) {
	s.mc.Finish()
}

func testCtx() context.Context {
	return context.Background()
}
