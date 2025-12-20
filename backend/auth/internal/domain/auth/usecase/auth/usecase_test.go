package auth

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/entity"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase/mocks"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/db"
	"github.com/m11ano/budget_planner/backend/auth/pkg/pgclient"
	"github.com/stretchr/testify/require"
)

type dependencies struct {
	ctx               context.Context
	cfg               config.Config
	dbMasterClient    db.MasterClient
	accountUCMock     *mocks.AccountUsecaseMock
	sessionUCMock     *mocks.SessionUsecaseMock
	accountUC         usecase.AccountUsecase
	sessionUC         usecase.SessionUsecase
	uc                *UsecaseImpl
	jwtAccessSecret   string
	jwtRefreshSecret  string
	accessLifetimeSec int
	refreshLifetimeHr int
}

func newDependencies(t *testing.T) *dependencies {
	t.Helper()

	s := &dependencies{}
	s.ctx = context.Background()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	s.jwtAccessSecret = "test-access-secret"
	s.jwtRefreshSecret = "test-refresh-secret"
	s.accessLifetimeSec = 600
	s.refreshLifetimeHr = 72

	s.cfg = config.Config{}
	s.cfg.Auth.JwtAccessSecret = s.jwtAccessSecret
	s.cfg.Auth.JwtRefreshSecret = s.jwtRefreshSecret
	s.cfg.Auth.AccessTokenLifetimeSec = s.accessLifetimeSec
	s.cfg.Auth.RefreshTokenLifetimeHrs = s.refreshLifetimeHr

	s.dbMasterClient = pgclient.NewMock()

	s.accountUCMock = mocks.NewAccountUsecaseMock(t)
	s.sessionUCMock = mocks.NewSessionUsecaseMock(t)

	s.accountUC = s.accountUCMock
	s.sessionUC = s.sessionUCMock

	s.uc = NewUsecaseImpl(
		logger,
		s.cfg,
		s.dbMasterClient,
		s.accountUC,
		s.sessionUC,
	)

	return s
}

func finishDependencies(t *testing.T, s *dependencies) {
	t.Helper()
	s.accountUCMock.MinimockFinish()
	s.sessionUCMock.MinimockFinish()
}

type tableTest struct {
	name string
	data map[string]interface{}
	run  func(t *testing.T, s *dependencies)
}

func mustAccountDTO(t *testing.T, email, password string, confirmed bool, blocked bool) *usecase.AccountDTO {
	t.Helper()

	acc, err := entity.NewAccount(email, password, true, confirmed)
	require.NoError(t, err)
	acc.IsBlocked = blocked

	return &usecase.AccountDTO{Account: acc}
}
