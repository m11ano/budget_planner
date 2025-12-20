package main

import (
	"github.com/m11ano/budget_planner/backend/gateway/internal/app"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/config"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/fxboot"
	"go.uber.org/fx"
)

// @title Budget API
// @version 1.0
// @description API документация для budget-app
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig("configs/base.yml", "configs/base.local.yml")

	appOptions := fxboot.BackendAppGetOptionsMap(app.IDBackend, cfg)

	app := fx.New(
		fxboot.OptionsMapToSlice(appOptions)...,
	)

	app.Run()
}
