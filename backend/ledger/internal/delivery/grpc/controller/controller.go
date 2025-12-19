package controller

import (
	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	budgetUC "github.com/m11ano/budget_planner/backend/ledger/internal/domain/budget/usecase"
	"github.com/m11ano/budget_planner/backend/ledger/pkg/backoff"
	desc "github.com/m11ano/budget_planner/backend/ledger/pkg/proto_pb/ledger_service"
	"google.golang.org/grpc"
)

type controller struct {
	desc.UnimplementedLedgerServer
	pkg          string
	cfg          config.Config
	backoff      *backoff.Controller
	budgetFacade *budgetUC.Facade
}

func Register(
	gRPCServer *grpc.Server,
	cfg config.Config,
	backoffCtrl *backoff.Controller,
	budgetFacade *budgetUC.Facade,
) {
	ctrl := &controller{
		pkg:          "grpc.Controller",
		cfg:          cfg,
		backoff:      backoffCtrl,
		budgetFacade: budgetFacade,
	}

	desc.RegisterLedgerServer(gRPCServer, ctrl)
}
