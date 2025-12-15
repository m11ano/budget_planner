package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/m11ano/budget_planner/backend/auth/internal/infra/loghandler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func getClientIPFromGRPCRequest(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			parts := strings.Split(xff[0], ",")
			ip := strings.TrimSpace(parts[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}

		if xri := md.Get("x-real-ip"); len(xri) > 0 {
			ip := strings.TrimSpace(xri[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			return host
		}
		return p.Addr.String()
	}

	return ""
}

type contextReqIPKey string

const (
	contextReqIPKeyValue contextReqIPKey = "x-request-ip"
)

func interceptorRequestIP() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestIP := getClientIPFromGRPCRequest(ctx)

		ctx = context.WithValue(ctx, contextReqIPKeyValue, requestIP)
		ctx = loghandler.SetContextData(ctx, "request.ip", requestIP)

		return handler(ctx, req)
	}
}

func GetClientIP(ctx context.Context) string {
	if v := ctx.Value(contextReqIPKeyValue); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}

	return ""
}
