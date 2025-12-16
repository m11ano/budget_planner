// Package v1 contains v1 http handlers
package v1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/config"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/backend/middleware"
)

type Groups struct {
	Prefix  string
	Default fiber.Router
}

const BackoffDefaultGroupID = "default"

// ProvideGroups - provide v1 group
func ProvideGroups(
	cfg config.Config,
	fiberApp *fiber.App,
	mdwr *middleware.Controller,
) *Groups {
	prefix := fmt.Sprintf("%s/v1", cfg.BackendApp.HTTP.Prefix)

	defaultGroup := fiberApp.Group(prefix, mdwr.Auth)

	return &Groups{
		Prefix:  prefix,
		Default: defaultGroup,
	}
}
