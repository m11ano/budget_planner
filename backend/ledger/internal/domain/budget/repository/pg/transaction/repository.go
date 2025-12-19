package transaction

import (
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/db"
)

type Repository struct {
	pkg      string
	logger   *slog.Logger
	pgClient db.MasterClient
	qb       squirrel.StatementBuilderType
}

func NewRepository(logger *slog.Logger, pgClient db.MasterClient) *Repository {
	return &Repository{
		pkg:      "Budget.repository.Transaction",
		logger:   logger,
		pgClient: pgClient,
		qb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
