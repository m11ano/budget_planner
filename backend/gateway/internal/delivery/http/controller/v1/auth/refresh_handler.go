package auth

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/httperrs"
	auth_servicev1 "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type RefreshHandlerIn struct {
	RefreshJWT string `json:"refreshJWT"`
}

// RefreshHandler - refresh tokens
// @Summary Refresh tokens
// @Tags auth
// @Produce  json
// @Accept  json
// @Param request body RefreshHandlerIn true "JSON"
// @Success 200 {object} TokensOutDTO
// @Failure 400 {object} middleware.ErrorJSON
// @Router /auth/refresh [post]
func (ctrl *Controller) RefreshHandler(c *fiber.Ctx) error {
	const op = "RefreshHandler"

	in := &RefreshHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	data, err := ctrl.authAdapter.Api().Refresh(c.Context(), &auth_servicev1.RefreshRequest{
		RefreshJwt: in.RefreshJWT,
	})
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := TokensOutDTO{
		AccessJWT:  data.Tokens.AccessJwt,
		RefreshJWT: data.Tokens.RefreshJwt,
	}

	return c.JSON(out)
}
