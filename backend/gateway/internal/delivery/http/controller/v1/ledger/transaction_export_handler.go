package ledger

import (
	"cloud.google.com/go/civil"
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

// TransactionExportHandler - export list transactions
// @Summary Export list transactions
// @Security BearerAuth
// @Tags ledger
// @Param date_from query string false "Фильтр по дате ОТ в формате 2025-01-30 (год-месяц-день)"
// @Param date_to query string false "Фильтр по дате ДО в формате 2025-01-30 (год-месяц-день)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Produce text/csv
// @Success 200 "CSV file"
// @Failure 400 {object} middleware.ErrorJSON
// @Router /ledger/transactions/export [get]
func (ctrl *Controller) TransactionExportHandler(c *fiber.Ctx) error {
	const op = "TransactionExportHandler"

	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	request := &desc.ListTransactionsRequest{
		Limit:  int32(limit),
		Offset: int64(offset),
	}

	filterDateFromStr := c.Query("date_from")
	if filterDateFromStr != "" {
		filterDateFrom, err := civil.ParseDate(filterDateFromStr)
		if err != nil {
			return appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid date_from"),
				"%s.%s", ctrl.pkg, op,
			)
		}

		request.FilterOccurredOnFrom = &desc.Date{
			Year:  int32(filterDateFrom.Year),
			Month: int32(filterDateFrom.Month),
			Day:   int32(filterDateFrom.Day),
		}
	}

	filterDateToStr := c.Query("date_to")
	if filterDateToStr != "" {
		filterDateTo, err := civil.ParseDate(filterDateToStr)
		if err != nil {
			return appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid date_to"),
				"%s.%s", ctrl.pkg, op,
			)
		}

		request.FilterOccurredOnTo = &desc.Date{
			Year:  int32(filterDateTo.Year),
			Month: int32(filterDateTo.Month),
			Day:   int32(filterDateTo.Day),
		}
	}

	data, err := ctrl.ledgerAdapter.Api().CSVExportTransactions(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Attachment("transactions.csv")
	c.Set("Cache-Control", "no-store")

	return c.Send(data.Data)
}
