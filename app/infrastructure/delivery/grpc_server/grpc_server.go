package grpc_server

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/timsolov/ms-users/app/infrastructure/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegisterServiceCallback func(s grpc.ServiceRegistrar)

func Run(ctx context.Context, log logger.Logger, addr string, cb RegisterServiceCallback, serviceImpl interface{}) chan error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		log.Fatalf("bind port for gRPC server on %s", addr)
	}

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return status.Errorf(codes.Unknown, "panic triggered: %v", p)
		}),
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
	)

	cb(s)

	errCh := make(chan error, 1)

	// Serve gRPC Server
	log.Infof("Serving gRPC server on http://%s", addr)
	go func() {
		errCh <- s.Serve(lis)
		log.Infof("Shutdown gRPC server.")
	}()

	return errCh
}
