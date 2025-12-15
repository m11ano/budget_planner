package account

import (
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/db"
)

type Repository struct {
	pkg      string
	logger   *slog.Logger
	pgClient db.MasterClient
	qb       squirrel.StatementBuilderType
}

func NewRepository(logger *slog.Logger, pgClient db.MasterClient) *Repository {
	return &Repository{
		pkg:      "Auth.repository.Account",
		logger:   logger,
		pgClient: pgClient,
		qb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
