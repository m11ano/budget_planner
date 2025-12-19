package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	"github.com/m11ano/budget_planner/backend/gateway/internal/infra/loghandler"
)

func (ctrl *Controller) Auth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	if accessToken == "" {
		return c.Next()
	}

	claims, err := auth.ParseAccessToken(accessToken, true, []byte(ctrl.cfg.Auth.JwtAccessSecret))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return c.Next()
		}
		return err
	}

	c.Locals(auth.ContextKeyAccessToken, accessToken)

	authData, err := auth.ClaimsToAuthData(claims)
	if err != nil {
		return err
	}

	c.Locals(auth.ContextKeyAuthData, authData)

	ctxData, ctxKey := loghandler.SetData(c.Context(), "request.account.id", claims.AccountId)
	c.Locals(ctxKey, ctxData)

	return c.Next()
}

func GetAuthData(c *fiber.Ctx) *auth.AuthData {
	return auth.GetAuthData(c.Context())
}
