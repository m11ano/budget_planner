package middleware

import "github.com/m11ano/budget_planner/backend/gateway/internal/app/config"

type Controller struct {
	cfg config.Config
}

func New(cfg config.Config) *Controller {
	return &Controller{
		cfg: cfg,
	}
}
