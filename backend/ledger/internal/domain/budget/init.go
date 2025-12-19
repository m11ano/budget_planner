package budget

import (
	"log/slog"

	"github.com/m11ano/budget_planner/backend/ledger/internal/app"
	"github.com/m11ano/budget_planner/backend/ledger/internal/app/fxboot/invoking"
)

func Init(_ app.ID, logger *slog.Logger) invoking.InvokeInit {
	return invoking.InvokeInit{}
}
