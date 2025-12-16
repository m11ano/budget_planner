package main

import (
	"github.com/m11ano/budget_planner/backend/gateway/internal/app"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/config"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/fxboot"
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
