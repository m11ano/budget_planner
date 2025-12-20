package ledger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

type BudgetGetHandlerOutput struct {
	Item     *BudgetOutput `json:"item"`
	HitCache bool          `json:"hitCache"`
}

func (ctrl *Controller) BudgetGetHandler(c *fiber.Ctx) error {
	const op = "BudgetGetHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"),
			"%s.%s", ctrl.pkg, op,
		)
	}

	request := &desc.GetBudgetRequest{
		Id: id.String(),
	}

	data, err := ctrl.ledgerAdapter.Api().GetBudget(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := BudgetGetHandlerOutput{
		Item:     NewBudgetOutput(data.Item),
		HitCache: data.HitCache,
	}

	return c.JSON(out)
}
