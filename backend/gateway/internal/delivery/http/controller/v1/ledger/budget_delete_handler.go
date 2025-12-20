package ledger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

// BudgetDeleteHandler - delete budget
// @Summary Delete budget
// @Security BearerAuth
// @Tags ledger
// @Param id path string true "ID"
// @Success 200
// @Failure 400 {object} middleware.ErrorJSON
// @Router /ledger/budgets/{id} [delete]
func (ctrl *Controller) BudgetDeleteHandler(c *fiber.Ctx) error {
	const op = "BudgetDeleteHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"),
			"%s.%s", ctrl.pkg, op,
		)
	}

	request := &desc.DeleteBudgetRequest{
		Id: id.String(),
	}

	_, err = ctrl.ledgerAdapter.Api().DeleteBudget(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	return nil
}
