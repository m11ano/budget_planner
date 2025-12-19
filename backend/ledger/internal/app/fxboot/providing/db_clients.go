package providing

import (
	"context"
	"log/slog"

	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/db"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/pgclient"
	"go.uber.org/fx"
)

func NewDBClients(masterDSN string, logQueries bool, logger *slog.Logger, shutdown fx.Shutdowner) db.MasterClient {
	master, err := pgclient.NewClient(
		context.Background(),
		"master",
		masterDSN,
		pgclient.NewClientOpts{
			Logger:     logger,
			LogQueries: logQueries,
		},
	)
	if err != nil {
		logger.Error("failed to create master client", slog.Any("error", err))
		// nolint
		_ = shutdown.Shutdown()
	}

	return master
}
