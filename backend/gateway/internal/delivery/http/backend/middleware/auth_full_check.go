package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) AuthFullCheck(c *fiber.Ctx) error {
	// authData := GetAuthData(c)
	// if authData == nil {
	// 	return appErrors.ErrUnauthorized
	// }

	// c.Locals(auth.ContextKeyAuthCheckRight, true)

	// isConfirmed, err := ctrl.authUC.IsSessionConfirmed(c.Context(), authData.SessionID)
	// if err != nil {
	// 	return err
	// }

	// if !isConfirmed {
	// 	return appErrors.ErrUnauthorized
	// }

	return c.Next()
}
