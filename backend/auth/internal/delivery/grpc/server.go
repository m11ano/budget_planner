package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/m11ano/budget_planner/backend/auth/internal/app/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func New(logger *slog.Logger, cfg config.Config) *grpc.Server {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.FinishCall,
			logging.PayloadReceived,
			logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		interceptorErrors(),
		recovery.UnaryServerInterceptor(recoveryOpts...),
		interceptorRequestIP(),
		interceptorRequestID(),
		interceptorAuth(cfg),
		interceptorValidate,
		logging.UnaryServerInterceptor(interceptorLogger(logger), loggingOpts...),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	reflection.Register(gRPCServer)

	return gRPCServer
}

func Start(grpcServer *grpc.Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %w", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}
