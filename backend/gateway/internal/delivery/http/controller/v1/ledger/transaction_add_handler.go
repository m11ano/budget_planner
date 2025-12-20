package ledger

import (
	"cloud.google.com/go/civil"
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/httperrs"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type TransactionAddHandlerInput struct {
	Amount      string     `json:"amount"`
	CategoryID  uint64     `json:"categoryID"`
	Description string     `json:"description"`
	IsIncome    bool       `json:"isIncome"`
	OccurredOn  civil.Date `json:"occurredOn" swaggertype:"string" example:"2025-12-20"`
}

type TransactionAddHandlerOutput struct {
	Item *TransactionOutput `json:"item"`
}

// TransactionAddHandler - add transaction
// @Summary Add transaction
// @Security BearerAuth
// @Tags ledger
// @Accept  json
// @Produce  json
// @Param request body TransactionAddHandlerInput true "JSON"
// @Success 200 {object} TransactionAddHandlerOutput
// @Failure 400 {object} middleware.ErrorJSON
// @Router /ledger/transactions [post]
func (ctrl *Controller) TransactionAddHandler(c *fiber.Ctx) error {
	const op = "TransactionAddHandler"

	in := &TransactionAddHandlerInput{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	request := &desc.AddTransactionRequest{
		Amount:   in.Amount,
		IsIncome: in.IsIncome,
		OccurredOn: &desc.Date{
			Year:  int32(in.OccurredOn.Year),
			Month: int32(in.OccurredOn.Month),
			Day:   int32(in.OccurredOn.Day),
		},
		CategoryId:  int64(in.CategoryID),
		Description: in.Description,
	}

	data, err := ctrl.ledgerAdapter.Api().AddTransaction(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := TransactionAddHandlerOutput{
		Item: NewTransactionOutput(data.Item),
	}

	return c.JSON(out)
}
