package auth

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/middleware"
	auth_servicev1 "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
)

// WhoIAmHandler - who i am
// @Summary Who i am
// @Security BearerAuth
// @Tags auth
// @Produce  json
// @Success 200 {object} AccountOutDTO
// @Failure 400 {object} middleware.ErrorJSON
// @Router /auth/whoiam [get]
func (ctrl *Controller) WhoIAmHandler(c *fiber.Ctx) error {
	const op = "WhoIAmHandler"

	authData := middleware.GetAuthData(c)
	if authData == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	resp, err := ctrl.authAdapter.Api().WhoIAm(c.Context(), &auth_servicev1.WhoIAmRequest{})
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := AccountOutDTO{
		ID:             resp.Account.Id,
		Email:          resp.Account.Email,
		ProfileName:    resp.Account.ProfileName,
		ProfileSurname: resp.Account.ProfileSurname,
	}

	return c.JSON(out)
}
