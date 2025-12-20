package ledger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

type TransactionGetHandlerOutput struct {
	Item *TransactionOutput `json:"item"`
}

func (ctrl *Controller) TransactionGetHandler(c *fiber.Ctx) error {
	const op = "TransactionGetHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(
			appErrors.ErrBadRequest.WithWrap(err).WithHints("invalid id"),
			"%s.%s", ctrl.pkg, op,
		)
	}

	request := &desc.GetTransactionRequest{
		Id: id.String(),
	}

	data, err := ctrl.ledgerAdapter.Api().GetTransaction(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := TransactionGetHandlerOutput{
		Item: NewTransactionOutput(data.Item),
	}

	return c.JSON(out)
}
