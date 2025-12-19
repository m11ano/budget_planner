package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/m11ano/budget_planner/backend/auth/pkg/auth"
	"github.com/m11ano/budget_planner/backend/ledger/internal/app/config"
	"github.com/m11ano/budget_planner/backend/ledger/internal/infra/loghandler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func interceptorAuth(cfg config.Config) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return handler(ctx, req)
		}

		authHeader := authHeaders[0]
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		if accessToken == "" {
			return handler(ctx, req)
		}

		claims, err := auth.ParseAccessToken(accessToken, true, []byte(cfg.Auth.JwtAccessSecret))
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) {
				return handler(ctx, req)
			}
			return nil, err
		}

		authData, err := auth.ClaimsToAuthData(claims)
		if err != nil {
			return nil, err
		}

		ctx = auth.SetAuthData(ctx, authData)

		ctx = loghandler.SetContextData(ctx, "request.account.id", claims.AccountId)

		return handler(ctx, req)
	}
}
