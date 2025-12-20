package session

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase/mocks"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/db"
	"github.com/m11ano/budget_planner/backend/auth/pkg/pgclient"
)

type dependencies struct {
	ctx            context.Context
	cfg            config.Config
	dbMasterClient db.MasterClient
	repo           usecase.SessionRepository
	repoMock       *mocks.SessionRepositoryMock
	uc             *UsecaseImpl
}

func newDependencies(t *testing.T) *dependencies {
	s := &dependencies{}
	s.ctx = context.Background()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	s.cfg = config.Config{}

	s.dbMasterClient = pgclient.NewMock()

	s.repoMock = mocks.NewSessionRepositoryMock(t)
	s.repo = s.repoMock

	s.uc = NewUsecaseImpl(
		logger,
		s.cfg,
		s.dbMasterClient,
		s.repo,
	)

	return s
}

func finishDependencies(t *testing.T, s *dependencies) {
	t.Helper()
	s.repoMock.MinimockFinish()
}

type tableTest struct {
	name string
	data map[string]interface{}
	run  func(t *testing.T, s *dependencies)
}
