package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
)

func (ctrl *Controller) AuthRequired(c *fiber.Ctx) error {
	authData := GetAuthData(c)
	if authData == nil {
		return appErrors.ErrUnauthorized
	}

	c.Locals(auth.ContextKeyAuthCheckRight, true)

	return c.Next()
}
