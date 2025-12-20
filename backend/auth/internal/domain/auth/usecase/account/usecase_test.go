package account

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase/mocks"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/db"
	"github.com/m11ano/budget_planner/backend/auth/pkg/pgclient"
)

type dependencies struct {
	ctx            context.Context
	cfg            config.Config
	dbMasterClient db.MasterClient

	accountRepo     usecase.AccountRepository
	accountRepoMock *mocks.AccountRepositoryMock

	sessionUC     usecase.SessionUsecase
	sessionUCMock *mocks.SessionUsecaseMock

	uc *UsecaseImpl
}

func newDependencies(t *testing.T) *dependencies {
	s := &dependencies{}
	s.ctx = context.Background()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	s.cfg = config.Config{}

	s.dbMasterClient = pgclient.NewMock()

	s.accountRepoMock = mocks.NewAccountRepositoryMock(t)
	s.accountRepo = s.accountRepoMock

	s.sessionUCMock = mocks.NewSessionUsecaseMock(t)
	s.sessionUC = s.sessionUCMock

	s.uc = NewUsecaseImpl(
		logger,
		s.cfg,
		s.dbMasterClient,
		s.accountRepo,
		s.sessionUC,
	)

	return s
}

func finishDependencies(t *testing.T, s *dependencies) {
	t.Helper()
	s.accountRepoMock.MinimockFinish()
	s.sessionUCMock.MinimockFinish()
}

type tableTest struct {
	name string
	data map[string]interface{}
	run  func(t *testing.T, s *dependencies)
}

func newAccount(id uuid.UUID, email string) *entity.Account {
	return &entity.Account{
		ID:        id,
		Email:     email,
		UpdatedAt: time.Now(),
	}
}
