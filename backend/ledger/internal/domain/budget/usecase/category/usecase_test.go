package category

import (
	"io"
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase/mocks"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/db"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/pgclient"
)

type dependencies struct {
	mc *minimock.Controller

	logger *slog.Logger
	cfg    config.Config

	dbMasterClient db.MasterClient

	categoryRepo *mocks.CategoryRepositoryMock

	uc *UsecaseImpl
}

func newDependencies(t *testing.T) *dependencies {
	t.Helper()

	mc := minimock.NewController(t)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := config.Config{}

	dbMasterClient := pgclient.NewMock()

	categoryRepo := mocks.NewCategoryRepositoryMock(mc)

	uc := NewUsecaseImpl(
		logger,
		cfg,
		dbMasterClient,
		categoryRepo,
	)

	return &dependencies{
		mc:             mc,
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		categoryRepo:   categoryRepo,
		uc:             uc,
	}
}

func finishDependencies(t *testing.T, d *dependencies) {
	t.Helper()
	require.NotNil(t, d)
}
