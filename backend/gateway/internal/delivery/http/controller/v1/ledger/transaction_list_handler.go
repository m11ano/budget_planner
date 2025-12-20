package ledger

import (
	"cloud.google.com/go/civil"
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

type TransactionListHandlerOutput struct {
	Items []*TransactionOutput `json:"items"`
	Total int64                `json:"total"`
}

// TransactionListHandler - list transactions
// @Summary List transactions
// @Security BearerAuth
// @Tags ledger
// @Param date_from query string false "Фильтр по дате ОТ в формате 2025-01-30 (год-месяц-день)"
// @Param date_to query string false "Фильтр по дате ДО в формате 2025-01-30 (год-месяц-день)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Produce json
// @Success 200 {object} TransactionListHandlerOutput
// @Failure 400 {object} middleware.ErrorJSON
// @Router /ledger/transactions [get]
func (ctrl *Controller) TransactionListHandler(c *fiber.Ctx) error {
	const op = "TransactionListHandler"

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

	data, err := ctrl.ledgerAdapter.Api().ListTransactions(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := TransactionListHandlerOutput{
		Items: make([]*TransactionOutput, 0, len(data.Items)),
		Total: data.Total,
	}

	for _, data := range data.Items {
		out.Items = append(out.Items, NewTransactionOutput(data))
	}

	return c.JSON(out)
}
