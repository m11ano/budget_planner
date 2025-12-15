package auth

import (
	"log/slog"

	"github.com/m11ano/budget_planner/backend/auth/internal/app"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/fxboot/invoking"
)

func Init(_ app.ID, logger *slog.Logger) invoking.InvokeInit {
	return invoking.InvokeInit{}
}
