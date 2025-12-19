package auth

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/middleware"
	auth_servicev1 "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
)

func (ctrl *Controller) LogoutHandler(c *fiber.Ctx) error {
	const op = "LogoutHandler"

	authData := middleware.GetAuthData(c)
	if authData == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	_, err := ctrl.authAdapter.Api().Logout(c.Context(), &auth_servicev1.LogoutRequest{})
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	return nil
}
