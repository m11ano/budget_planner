package auth

import (
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
	"google.golang.org/grpc"
)

type Adapter interface {
	Api() desc.AuthClient
	CC() *grpc.ClientConn
}
