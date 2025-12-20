package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
)

type ErrorJSON struct {
	Code     int            `json:"code"`
	TextCode string         `json:"textCode"`
	Hints    []string       `json:"hints"`
	Details  map[string]any `json:"details"`
}

func ErrorHandler(appTitle string, logger *slog.Logger) func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		code := 500
		jsonRes := ErrorJSON{
			TextCode: "INTERNAL_ERROR",
			Hints:    []string{},
			Details:  map[string]any{},
		}

		if appError, ok := appErrors.ExtractError(err); ok {
			code = int(appError.Meta().Code)
			jsonRes.TextCode = appError.Meta().TextCode
			jsonRes.Hints = appError.Hints()
			jsonRes.Details = appError.Details(false)
		} else {
			switch errTyped := err.(type) {
			case *fiber.Error:
				code = errTyped.Code
				switch {
				case code == 405:
					jsonRes.TextCode = "METHOD_NOT_ALLOWED"
				case code >= 400 && code < 500:
					jsonRes.TextCode = "BAD_REQUEST"
				}
				jsonRes.Hints = []string{errTyped.Message}
			default:
			}
		}

		jsonRes.Code = code

		return c.Status(code).JSON(jsonRes)
	}
}
