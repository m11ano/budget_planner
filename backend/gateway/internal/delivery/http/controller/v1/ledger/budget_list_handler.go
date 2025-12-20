package ledger

import (
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

type BudgetListHandlerOutput struct {
	Items    []*BudgetOutput `json:"items"`
	Total    int64           `json:"total"`
	HitCache bool            `json:"hit_cache"`
}

func (ctrl *Controller) BudgetListHandler(c *fiber.Ctx) error {
	const op = "BudgetListHandler"

	limit := c.QueryInt("limit", 100)
	offset := c.QueryInt("offset", 0)

	request := &desc.ListBudgetsRequest{
		Limit:  int32(limit),
		Offset: int64(offset),
	}

	filterDateFromStr := c.Query("period_from")
	if filterDateFromStr != "" {
		filterDateFrom, err := civil.ParseDate(fmt.Sprintf("%s-01", filterDateFromStr))
		if err != nil {
			return appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid period_from"),
				"%s.%s", ctrl.pkg, op,
			)
		}

		request.FilterPeriodFrom = &desc.DateMonth{
			Year:  int32(filterDateFrom.Year),
			Month: int32(filterDateFrom.Month),
		}
	}

	filterDateToStr := c.Query("period_to")
	if filterDateToStr != "" {
		filterDateTo, err := civil.ParseDate(fmt.Sprintf("%s-01", filterDateToStr))
		if err != nil {
			return appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid period_to"),
				"%s.%s", ctrl.pkg, op,
			)
		}

		request.FilterPeriodTo = &desc.DateMonth{
			Year:  int32(filterDateTo.Year),
			Month: int32(filterDateTo.Month),
		}
	}

	data, err := ctrl.ledgerAdapter.Api().ListBudgets(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := BudgetListHandlerOutput{
		Items:    make([]*BudgetOutput, 0, len(data.Items)),
		Total:    data.Total,
		HitCache: data.HitCache,
	}

	for _, data := range data.Items {
		out.Items = append(out.Items, NewBudgetOutput(data))
	}

	return c.JSON(out)
}
