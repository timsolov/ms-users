package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net"
	"os"
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

	// TODO: add methods for LoggerV2 to LogrusLogger
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := conf.New()

	log2 := logger.NewLogrusLogger(cfg.LOG().LogLevel, cfg.LOG().LogJson, "", false)

	d, err := postgres.New(cfg.DB().DSN(), 5, 5, 5*time.Minute)
	if err != nil {
		log2.Fatalf("connect to db: %v", err)
	}

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	server.RegisterUserServiceServer(s, server.New(usecase.New(d)))

	// Serve gRPC Server
	log.Info("Serving gRPC on http://", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run(ctx, "0.0.0.0:11000", addr, []gateway.RegisterServiceHandlerFunc{server.RegisterUserServiceHandler})
	log.Fatalln(err)
}
