package auth

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/backend/middleware"
	auth_servicev1 "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
)

func (ctrl *Controller) WhoIAmHandler(c *fiber.Ctx) error {
	authData := middleware.GetAuthData(c)
	if authData == nil {
		return appErrors.ErrUnauthorized
	}

	resp, err := ctrl.authAdapter.Api().GetAccountByID(c.Context(), &auth_servicev1.GetAccountByIDRequest{
		AccountId: authData.AccountID.String(),
	})
	if err != nil {
		return err
	}

	out := AccoutOutDTO{
		ID:             resp.Account.Id,
		Email:          resp.Account.Email,
		ProfileName:    resp.Account.ProfileName,
		ProfileSurname: resp.Account.ProfileSurname,
	}

	return c.JSON(out)
}
