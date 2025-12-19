package controller

import (
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/middleware"
	v1 "github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/v1"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/v1/auth"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/v1/tests"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
	"go.uber.org/fx"
)

// FxModule - fx module
var FxModule = fx.Options(
	fx.Provide(validation.New),
	fx.Options(
		fx.Provide(middleware.New),
		fx.Provide(v1.ProvideGroups),
		tests.FxModule,
		auth.FxModule,
	),
)
