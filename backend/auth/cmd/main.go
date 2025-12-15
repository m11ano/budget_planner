package main

import (
	"github.com/m11ano/budget_planner/backend/auth/internal/app"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/fxboot"
	"go.uber.org/fx"
)

func main() {
	cfg := config.LoadConfig("configs/base.yml", "configs/base.local.yml")

	appOptions := fxboot.BackendAppGetOptionsMap(app.IDBackend, cfg)

	app := fx.New(
		fxboot.OptionsMapToSlice(appOptions)...,
	)

	app.Run()
}
