package auth

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	authAdapter "github.com/m11ano/budget_planner/backend/gateway/internal/adapter/auth"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/config"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/middleware"
	v1 "github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/v1"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type Controller struct {
	pkg         string
	vldtr       *validator.Validate
	cfg         config.Config
	authAdapter authAdapter.Adapter
}

func NewController(
	cfg config.Config,
	authAdapter authAdapter.Adapter,
) *Controller {
	controller := &Controller{
		pkg:         "httpController.Auth",
		vldtr:       validation.New(),
		cfg:         cfg,
		authAdapter: authAdapter,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, mdwr *middleware.Controller) {
	const url = "auth"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.Post("/register", ctrl.RegisterHandler)

	routeGroup.Post("/login", ctrl.LoginHandler)

	routeGroup.Post("/refresh", ctrl.RefreshHandler)

	routeGroup.Post("/logout", ctrl.LogoutHandler)

	routeGroup.Get("/whoiam", ctrl.WhoIAmHandler)
}
