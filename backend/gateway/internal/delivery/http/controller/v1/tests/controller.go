package tests

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/config"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/middleware"
	v1 "github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/v1"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type Controller struct {
	pkg   string
	vldtr *validator.Validate
	cfg   config.Config
}

func NewController(
	cfg config.Config,
) *Controller {
	controller := &Controller{
		pkg:   "httpController.Tests",
		vldtr: validation.New(),
		cfg:   cfg,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, mdwr *middleware.Controller) {
	const url = "tests"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.All("/ping", ctrl.PingHandler)
}
