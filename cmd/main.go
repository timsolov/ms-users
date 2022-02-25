package main

import (
	"context"
	"flag"
	"net"
	"time"

	"github.com/timsolov/ms-users/internal/client/db/postgres"
	"github.com/timsolov/ms-users/internal/common/gateway"
	"github.com/timsolov/ms-users/internal/common/logger"
	"github.com/timsolov/ms-users/internal/conf"
	"github.com/timsolov/ms-users/internal/server"
	"github.com/timsolov/ms-users/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	flag.Parse()

	cfg := conf.New()
	log := logger.NewLogrusLogger(cfg.LOG().LogLevel, cfg.LOG().LogJson, "", false)
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	grpclog.SetLoggerV2(log.(grpclog.LoggerV2))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d, err := postgres.New(cfg.DB().DSN(), 5, 5, 5*time.Minute)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Bind port for gRPC server on http://%s", addr)
	}

	s := grpc.NewServer()
	server.RegisterUserServiceServer(s, server.New(usecase.New(d)))

	// Serve gRPC Server
	log.Infof("Serving gRPC on http://%s", addr)
	go func() {
		log.Fatalf("%s", s.Serve(lis))
	}()

	err = gateway.Run(ctx, "0.0.0.0:11000", addr, []gateway.RegisterServiceHandlerFunc{server.RegisterUserServiceHandler})
	log.Fatalf("Run gateway: %s", err)
}
