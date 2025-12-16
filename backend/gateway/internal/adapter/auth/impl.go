package auth

import (
	"log/slog"
	"time"

	appErrors "github.com/m11ano/budget_planner/backend/gateway/internal/app/errors"
	grpcClient "github.com/m11ano/budget_planner/backend/gateway/internal/client/grpc"
	desc "github.com/m11ano/budget_planner/backend/gateway/pkg/proto_pb/auth_service"
	"google.golang.org/grpc"
)

type AdapterImpl struct {
	api    desc.AuthClient
	cc     *grpc.ClientConn
	logger *slog.Logger
}

func NewAdapterImpl(
	addr string,
	retriesCount int,
	timeout time.Duration,
	logger *slog.Logger,
) (*AdapterImpl, error) {
	const op = "NewAdapterImpl"

	cfg := grpcClient.Config{
		Addr:         addr,
		RetriesCount: retriesCount,
		Timeout:      timeout,
	}

	cc, err := grpcClient.NewClientConn(cfg, logger)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s", op)
	}

	adapter := &AdapterImpl{
		api:    desc.NewAuthClient(cc),
		logger: logger,
		cc:     cc,
	}

	return adapter, nil
}

func (a *AdapterImpl) Api() desc.AuthClient {
	return a.api
}

func (a *AdapterImpl) CC() *grpc.ClientConn {
	return a.cc
}
