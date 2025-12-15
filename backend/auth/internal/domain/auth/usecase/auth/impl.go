package auth

import (
	"log/slog"

	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/db"

	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	dbMasterClient db.MasterClient
	accountUC      usecase.AccountUsecase
	sessionUC      usecase.SessionUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	accountUC usecase.AccountUsecase,
	sessionUC usecase.SessionUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "Auth.usecase.Auth",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		accountUC:      accountUC,
		sessionUC:      sessionUC,
	}

	return uc
}
