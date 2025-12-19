package auth

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/httperrs"
	auth_servicev1 "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type RegisterHandlerIn struct {
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"min=0"`
	ProfileName    string `json:"profileName" validate:"required"`
	ProfileSurname string `json:"profileSurname" validate:"required"`
}

func (ctrl *Controller) RegisterHandler(c *fiber.Ctx) error {
	const op = "RegisterHandler"

	in := &RegisterHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	_, err := ctrl.authAdapter.Api().Register(c.Context(), &auth_servicev1.RegisterRequest{
		Email:          in.Email,
		Password:       in.Password,
		ProfileName:    in.ProfileName,
		ProfileSurname: in.ProfileSurname,
	})
	if err != nil {
		return appErrors.Chainf(appErrors.FromGRPCError(err), "%s.%s", ctrl.pkg, op)
	}

	return nil
}
