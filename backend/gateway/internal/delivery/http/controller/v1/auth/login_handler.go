package auth

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/httperrs"
	auth_servicev1 "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type LoginHandlerIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=0"`
}

type LoginHandlerOut struct {
	Tokens  TokensOutDTO  `json:"tokens"`
	Account AccountOutDTO `json:"account"`
}

// LoginHandler - Login
// @Summary Login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body LoginHandlerIn true "JSON"
// @Success 200 {object} LoginHandlerOut
// @Failure 400 {object} middleware.ErrorJSON
// @Router /auth/login [post]
func (ctrl *Controller) LoginHandler(c *fiber.Ctx) error {
	const op = "LoginHandler"

	in := &LoginHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	data, err := ctrl.authAdapter.Api().Login(c.Context(), &auth_servicev1.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	out := LoginHandlerOut{
		Tokens: TokensOutDTO{
			AccessJWT:  data.Tokens.AccessJwt,
			RefreshJWT: data.Tokens.RefreshJwt,
		},
		Account: AccountOutDTO{
			ID:             data.Account.Id,
			Email:          data.Account.Email,
			ProfileName:    data.Account.ProfileName,
			ProfileSurname: data.Account.ProfileSurname,
		},
	}

	return c.JSON(out)
}
