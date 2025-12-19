package transaction

import (
	"log/slog"

	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/db"
)

type UsecaseImpl struct {
	pkg             string
	logger          *slog.Logger
	cfg             config.Config
	dbMasterClient  db.MasterClient
	transactionRepo usecase.TransactionRepository
	categoryRepo    usecase.CategoryRepository
	budgetRepo      usecase.BudgetRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	transactionRepo usecase.TransactionRepository,
	categoryRepo usecase.CategoryRepository,
	budgetRepo usecase.BudgetRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:             "Budget.Usecase.Transaction",
		logger:          logger,
		cfg:             cfg,
		dbMasterClient:  dbMasterClient,
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
		budgetRepo:      budgetRepo,
	}
	return uc
}
