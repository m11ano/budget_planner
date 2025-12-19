package tests

import (
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) PingHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}
