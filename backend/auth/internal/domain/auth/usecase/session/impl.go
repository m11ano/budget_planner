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
	repo           usecase.SessionRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	repo usecase.SessionRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "TenantUser.usecase.Session",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		repo:           repo,
	}

	return uc
}
