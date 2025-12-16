package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/gateway/internal/infra/loghandler"
)

func RequestID() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestID, err := uuid.Parse(c.Get("X-Request-ID"))
		if err != nil {
			requestID = uuid.New()
		}

		c.Locals("requestID", requestID)

		ctxData, ctxKey := loghandler.SetData(c.Context(), "request.id", requestID)
		c.Locals(ctxKey, ctxData)

		return c.Next()
	}
}
