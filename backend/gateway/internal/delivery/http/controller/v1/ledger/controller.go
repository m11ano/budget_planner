package ledger

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	ledgerAdapter "github.com/m11ano/budget_planner/backend/gateway/internal/adapter/ledger"
	"github.com/m11ano/budget_planner/backend/gateway/internal/app/config"
	"github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/middleware"
	v1 "github.com/m11ano/budget_planner/backend/gateway/internal/delivery/http/controller/v1"
	"github.com/m11ano/budget_planner/backend/gateway/pkg/validation"
)

type Controller struct {
	pkg           string
	vldtr         *validator.Validate
	cfg           config.Config
	ledgerAdapter ledgerAdapter.Adapter
}

func NewController(
	cfg config.Config,
	ledgerAdapter ledgerAdapter.Adapter,
) *Controller {
	controller := &Controller{
		pkg:           "httpController.Ledger",
		vldtr:         validation.New(),
		cfg:           cfg,
		ledgerAdapter: ledgerAdapter,
	}
	return controller
}

func RegisterRoutes(groups *v1.Groups, ctrl *Controller, mdwr *middleware.Controller) {
	const url = "ledger"

	routeGroup := groups.Default.Group(fmt.Sprintf("/%s", url))

	routeGroup.Get("/transactions/export", ctrl.TransactionExportHandler)

	routeGroup.Post("/transactions/import", ctrl.TransactionImportHandler)
}
