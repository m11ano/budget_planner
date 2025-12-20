package ledger

import (
	"cloud.google.com/go/civil"
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

type ReportListHandlerOutput struct {
	Reports  []*ReportOutput `json:"reports"`
	HitCache bool            `json:"hitCache"`
}

func (ctrl *Controller) ReportListHandler(c *fiber.Ctx) error {
	const op = "ReportListHandler"

	request := &desc.ListReportsRequest{}

	filterDateFromStr := c.Query("date_from")
	if filterDateFromStr != "" {
		filterDateFrom, err := civil.ParseDate(filterDateFromStr)
		if err != nil {
			return appErrors.Chainf(
				appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid date_from"),
				"%s.%s", ctrl.pkg, op,
			)
		}

		request.DateFrom = &desc.Date{
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

		request.DateTo = &desc.Date{
			Year:  int32(filterDateTo.Year),
			Month: int32(filterDateTo.Month),
			Day:   int32(filterDateTo.Day),
		}
	}

	data, err := ctrl.ledgerAdapter.Api().ListReports(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := ReportListHandlerOutput{
		Reports:  make([]*ReportOutput, 0, len(data.Reports)),
		HitCache: data.HitCache,
	}

	for _, data := range data.Reports {
		out.Reports = append(out.Reports, NewReportOutput(data))
	}

	return c.JSON(out)
}
