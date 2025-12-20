package controller

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/m11ano/budget_planner/backend/gateway/docs"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/middleware"
)

const defaultBodyLimit = 10 * 1024 * 1024

type HTTPConfig struct {
	AppTitle         string
	UnderProxy       bool
	UseLogger        bool
	BodyLimit        int
	CorsAllowOrigins []string
	ServerIPs        []string
}

func NewHTTPFiber(httpCfg HTTPConfig, logger *slog.Logger) *fiber.App {
	if httpCfg.BodyLimit == -1 {
		httpCfg.BodyLimit = defaultBodyLimit
	}

	fiberCfg := fiber.Config{
		ErrorHandler: middleware.ErrorHandler(httpCfg.AppTitle, logger),
		BodyLimit:    httpCfg.BodyLimit,
	}

	if httpCfg.UnderProxy {
		fiberCfg.ProxyHeader = fiber.HeaderXForwardedFor
	}

	app := fiber.New(fiberCfg)

	app.Use(middleware.Recovery(logger))
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestIP(httpCfg.ServerIPs))

	if len(httpCfg.CorsAllowOrigins) > 0 {
		app.Use(middleware.Cors(httpCfg.CorsAllowOrigins))
	}

	if httpCfg.UseLogger {
		app.Use(middleware.Logger(logger))
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
