package ledger

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

func (ctrl *Controller) TransactionExportHandler(c *fiber.Ctx) error {
	const op = "TransactionExportHandler"

	request := &desc.ListTransactionsRequest{}

	data, err := ctrl.ledgerAdapter.Api().CSVExportTransactions(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="transactions.csv"`)
	c.Set("Cache-Control", "no-store")

	return c.Send(data.Data)
}
