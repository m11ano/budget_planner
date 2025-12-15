package account

import (
	"log/slog"

	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"github.com/m11ano/budget_planner/backend/auth/internal/domain/auth/usecase"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/db"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	dbMasterClient db.MasterClient
	accountRepo    usecase.AccountRepository
	sessionUC      usecase.SessionUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	accountRepo usecase.AccountRepository,
	sessionUC usecase.SessionUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "Auth.Usecase.Account",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		accountRepo:    accountRepo,
		sessionUC:      sessionUC,
	}
	return uc
}
