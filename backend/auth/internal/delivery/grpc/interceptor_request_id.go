package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/m11ano/budget_planner/backend/auth/internal/infra/loghandler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextReqIDKey string

const (
	requestIDKey         contextReqIDKey = "x-request-id"
	metadataRequestIDKey string          = "x-request-id"
)

func interceptorRequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var reqID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(metadataRequestIDKey); len(vals) > 0 {
				reqID = vals[0]
			}
		}

		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, requestIDKey, reqID)
		ctx = loghandler.SetContextData(ctx, "request.id", reqID)

		return handler(ctx, req)
	}
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(requestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}

	return ""
}
