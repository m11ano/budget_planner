package ledger

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
)

type CategoryListHandlerOutput struct {
	Items []*CategoryOutput `json:"items"`
}

func (ctrl *Controller) CategoryListHandler(c *fiber.Ctx) error {
	const op = "CategoryListHandler"

	request := &desc.ListCategoriesRequest{}

	data, err := ctrl.ledgerAdapter.Api().ListCategories(c.Context(), request)
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := CategoryListHandlerOutput{
		Items: make([]*CategoryOutput, 0, len(data.Items)),
	}

	for _, data := range data.Items {
		out.Items = append(out.Items, NewCategoryOutput(data))
	}

	return c.JSON(out)
}
