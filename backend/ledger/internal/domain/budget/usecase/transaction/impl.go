package transaction

import (
	"log/slog"

	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	"github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/db"
	"golang.org/x/sync/singleflight"
)

type UsecaseImpl struct {
	pkg                  string
	logger               *slog.Logger
	cfg                  config.Config
	dbMasterClient       db.MasterClient
	sfGroup              singleflight.Group
	transactionRepo      usecase.TransactionRepository
	transactionCacheRepo usecase.TransactionCacheRepository
	transactionCSVRepo   usecase.TransactionCSVRepository
	categoryRepo         usecase.CategoryRepository
	budgetRepo           usecase.BudgetRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	transactionRepo usecase.TransactionRepository,
	transactionCacheRepo usecase.TransactionCacheRepository,
	transactionCSVRepo usecase.TransactionCSVRepository,
	categoryRepo usecase.CategoryRepository,
	budgetRepo usecase.BudgetRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:                  "Budget.Usecase.Transaction",
		logger:               logger,
		cfg:                  cfg,
		dbMasterClient:       dbMasterClient,
		transactionRepo:      transactionRepo,
		transactionCacheRepo: transactionCacheRepo,
		transactionCSVRepo:   transactionCSVRepo,
		categoryRepo:         categoryRepo,
		budgetRepo:           budgetRepo,
	}
	return uc
}
