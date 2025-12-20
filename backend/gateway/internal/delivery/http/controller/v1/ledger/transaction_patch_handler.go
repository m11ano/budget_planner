package ledger

import (
	"cloud.google.com/go/civil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/httperrs"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
	"github.com/samber/lo"
)

type TransactionPatchHandlerInput struct {
	Amount      *string     `json:"amount"`
	CategoryID  *uint64     `json:"categoryID"`
	Description *string     `json:"description"`
	OccurredOn  *civil.Date `json:"occurredOn"`
}

type TransactionPatchHandlerOutput struct {
	Item *TransactionOutput `json:"item"`
}

func (ctrl *Controller) TransactionPatchHandler(c *fiber.Ctx) error {
	const op = "TransactionPatchHandler"

	in := &TransactionPatchHandlerInput{}

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

	request := &desc.PatchTransactionRequest{
		Id:          id.String(),
		Amount:      in.Amount,
		Description: in.Description,
	}

	if in.CategoryID != nil {
		request.CategoryId = lo.ToPtr(int64(*in.CategoryID))
	}

	if in.OccurredOn != nil {
		request.OccurredOn = &desc.Date{
			Year:  int32(in.OccurredOn.Year),
			Month: int32(in.OccurredOn.Month),
			Day:   int32(in.OccurredOn.Day),
		}
	}

	data, err := ctrl.ledgerAdapter.Api().PatchTransaction(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := TransactionPatchHandlerOutput{
		Item: NewTransactionOutput(data.Item),
	}

	return c.JSON(out)
}
