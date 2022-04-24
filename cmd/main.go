package main

import (
	"context"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	cfg := conf.New()

	log := logger.NewLogrusLogger(cfg.LOG.Level, cfg.LOG.Json, cfg.LOG.TimeFormat, false)
	grpclog.SetLoggerV2(log)

	log.Infof("application started")
	defer log.Infof("application finished")

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(done)

	go func() {
		<-done
		log.Errorf("SIGTERM detected")
		cancel()
	}()

	// create connetion to PostgreSQL
	d, err := postgres.New(
		ctx,
		cfg.DB.DSN(),
		postgres.SetMaxConns(cfg.DB.OpenLimit, cfg.DB.IdleLimit),
		postgres.SetConnsMaxLifeTime(cfg.DB.ConnLife, 0),
		postgres.SetLogger(log),
		postgres.SetReconnectTimeout(cfg.DB.ReconnectTimeout),
	)
	if err != nil {
		log.Errorf("connect to db: %v", err)
		return
	}

	// migrate db if need
	ParseParams(log, d)

	grpcErr := grpc_server.Run(
		ctx,
		log,
		cfg.GRPC.Addr(), // listen incoming host:port for gRPC server
		func(s grpc.ServiceRegistrar) {
			pb.RegisterUserServiceServer(s, web.New(log, d))
		},
	)

	grpcGwErr := grpc_gateway.Run(
		ctx,
		log,
		cfg.HTTP.Addr(), // listen incoming host:port for rest api
		cfg.GRPC.Addr(), // connect to gRPC server host:port
		[]grpc_gateway.RegisterServiceHandlerFunc{
			pb.RegisterUserServiceHandler,
		},
	)

	select {
	case <-ctx.Done():
	case err := <-grpcErr:
		log.Errorf("gRPC server error: %s", err)
	case err := <-grpcGwErr:
		log.Errorf("gRPC gateway error: %s", err)
	}

	cancel()

	time.Sleep(1 * time.Second)
}
