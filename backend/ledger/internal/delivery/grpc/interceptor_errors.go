package grpc

import (
	"context"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"google.golang.org/grpc"
)

func interceptorErrors() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		data, err := handler(ctx, req)

		return data, appErrors.ToGrpcStatus(err)
	}
}
