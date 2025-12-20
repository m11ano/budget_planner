package transaction

import (
	"log/slog"
)

type Repository struct {
	pkg    string
	logger *slog.Logger
}

func NewRepository(logger *slog.Logger) *Repository {
	return &Repository{
		pkg:    "Budget.repository.CSVTransaction",
		logger: logger,
	}
}
