package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/timsolov/ms-users/app/conf"
	"github.com/timsolov/ms-users/app/infrastructure/delivery/grpc_gateway"
	"github.com/timsolov/ms-users/app/infrastructure/delivery/grpc_server"
	"github.com/timsolov/ms-users/app/infrastructure/delivery/web"
	"github.com/timsolov/ms-users/app/infrastructure/logger"
	"github.com/timsolov/ms-users/app/infrastructure/repository/postgres"
	"github.com/timsolov/ms-users/app/pb"
	"github.com/timsolov/ms-users/app/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	flag.Parse()

	cfg := conf.New()

	log := logger.NewLogrusLogger(cfg.LOG().Level, cfg.LOG().Json, "", false)
	grpclog.SetLoggerV2(log.(grpclog.LoggerV2))

	log.Infof("application started")
	defer log.Infof("application finished")

	ctx := context.Background()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	d, err := postgres.New(
		timeoutCtx,
		cfg.DB().DSN(),
		cfg.DB().OpenLimit,
		cfg.DB().IdleLimit,
		cfg.DB().ConnLife,
	)
	if err != nil {
		log.Errorf("connect to db: %v", err)
		return
	}
	cancel()

	grpcCtx, grpcCancel := context.WithCancel(ctx)
	grpcErr := grpc_server.Run(
		grpcCtx,
		log,
		cfg.GRPC().Addr(), // listen incoming port for gRPC
		func(s grpc.ServiceRegistrar) {
			pb.RegisterUserServiceServer(s, web.New(usecase.New(d)))
		},
		web.New(usecase.New(d)),
	)

	grpcGwCtx, grpcGwCancel := context.WithCancel(ctx)
	grpcGwErr := grpc_gateway.Run(
		grpcGwCtx,
		log,
		cfg.HTTP().Addr(), // listen incoming port for rest api
		cfg.GRPC().Addr(), // addr for connection to gRPC server
		[]grpc_gateway.RegisterServiceHandlerFunc{
			pb.RegisterUserServiceHandler,
		},
	)

	select {
	case <-done:
		log.Infof("SIGTERM detected")
	case err := <-grpcErr:
		log.Errorf("gRPC server error: %s", err)
	case err := <-grpcGwErr:
		log.Errorf("gRPC gateway error: %s", err)
	}

	grpcCancel()
	grpcGwCancel()

	time.Sleep(1 * time.Second)
}
