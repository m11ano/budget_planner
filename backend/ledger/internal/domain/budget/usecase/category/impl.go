package category

import (
	"log/slog"

	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/db"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	dbMasterClient db.MasterClient
	categoryRepo   usecase.CategoryRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	categoryRepo usecase.CategoryRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "Budget.Usecase.Category",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		categoryRepo:   categoryRepo,
	}
	return uc
}
