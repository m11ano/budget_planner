package ledger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/httperrs"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
	"github.com/samber/lo"
)

type BudgetPatchHandlerInput struct {
	Amount     *string                      `json:"amount"`
	CategoryID *uint64                      `json:"categoryID"`
	Period     *BudgetAddHandlerInputPeriod `json:"period"`
}

type BudgetPatchHandlerOutput struct {
	Item *BudgetOutput `json:"item"`
}

func (ctrl *Controller) BudgetPatchHandler(c *fiber.Ctx) error {
	const op = "BudgetPatchHandler"

	in := &BudgetPatchHandlerInput{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"),
			"%s.%s", ctrl.pkg, op,
		)
	}

	request := &desc.PatchBudgetRequest{
		Id:     id.String(),
		Amount: in.Amount,
	}

	if in.CategoryID != nil {
		request.CategoryId = lo.ToPtr(int64(*in.CategoryID))
	}

	if in.Period != nil {
		request.Period = &desc.DateMonth{
			Year:  int32(in.Period.Year),
			Month: int32(in.Period.Month),
		}
	}

	data, err := ctrl.ledgerAdapter.Api().PatchBudget(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := BudgetPatchHandlerOutput{
		Item: NewBudgetOutput(data.Item),
	}

	return c.JSON(out)
}
