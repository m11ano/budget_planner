package ledger

import (
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/ledger_service"
	"google.golang.org/grpc"
)

type Adapter interface {
	Api() desc.LedgerClient
	CC() *grpc.ClientConn
}
