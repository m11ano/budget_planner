package grpc

import (
	"context"

	appErrors "github.com/m11ano/budget_planner/backend/ledger/internal/app/errors"
	"google.golang.org/grpc"
)

func interceptorValidate(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	if reqV, ok := req.(interface{ ValidateAll() error }); ok {
		if err := reqV.ValidateAll(); err != nil {
			return nil, appErrors.ErrBadRequest.WithWrap(err).WithHints(err.Error())
		}
	}
	return handler(ctx, req)
}
