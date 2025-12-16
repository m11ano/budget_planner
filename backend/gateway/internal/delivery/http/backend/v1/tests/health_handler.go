package tests

import (
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) HealthHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}
