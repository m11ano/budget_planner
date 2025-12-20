package ledger

import (
	"io"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

// TransactionImportHandler - import transactions from CSV
// @Summary Import transactions from CSV
// @Description Upload CSV file with transactions for import
// @Security BearerAuth
// @Tags ledger
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "CSV file with transactions"
// @Success 200 "Import started / completed successfully"
// @Failure 400 {object} middleware.ErrorJSON
// @Failure 500 {object} middleware.ErrorJSON
// @Router /ledger/transactions/import [post]
func (ctrl *Controller) TransactionImportHandler(c *fiber.Ctx) error {
	const op = "TransactionImportHandler"

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrBadRequest.WithWrap(err).WithHints("form field `file` is required"),
			"%s.%s", ctrl.pkg, op,
		)
	}

	f, err := fileHeader.Open()
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrInternal.WithWrap(err),
			"%s.%s", ctrl.pkg, op,
		)
	}
	// nolint
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrInternal.WithWrap(err),
			"%s.%s", ctrl.pkg, op,
		)
	}

	request := &desc.CSVImportTransactionsRequest{
		Data: fileData,
	}

	_, err = ctrl.ledgerAdapter.Api().CSVImportTransactions(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	return nil
}
