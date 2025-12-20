package ledger

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/httperrs"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type BudgetAddHandlerInputPeriod struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type BudgetAddHandlerInput struct {
	Amount     string                      `json:"amount"`
	CategoryID uint64                      `json:"categoryID"`
	Period     BudgetAddHandlerInputPeriod `json:"period"`
}

type BudgetAddHandlerOutput struct {
	Item *BudgetOutput `json:"item"`
}

func (ctrl *Controller) BudgetAddHandler(c *fiber.Ctx) error {
	const op = "BudgetAddHandler"

	in := &BudgetAddHandlerInput{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	request := &desc.AddBudgetRequest{
		Amount: in.Amount,
		Period: &desc.DateMonth{
			Year:  int32(in.Period.Year),
			Month: int32(in.Period.Month),
		},
		CategoryId: int64(in.CategoryID),
	}

	data, err := ctrl.ledgerAdapter.Api().AddBudget(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := BudgetAddHandlerOutput{
		Item: NewBudgetOutput(data.Item),
	}

	return c.JSON(out)
}
