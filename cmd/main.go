package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ms-users/app/common/logger"
	"ms-users/app/conf"
	"ms-users/app/infrastructure/delivery/cli"
	"ms-users/app/infrastructure/delivery/grpc_gateway"
	"ms-users/app/infrastructure/delivery/grpc_server"
	"ms-users/app/infrastructure/delivery/web"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/infrastructure/repository/postgres"
	"ms-users/app/usecase/auth_emailpass"
	"ms-users/app/usecase/confirm"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"
	"ms-users/app/usecase/retry_confirm"
	"ms-users/app/usecase/whoami"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	cfg := conf.New()

	log := logger.NewLogrusLogger(cfg.LOG.Level, cfg.LOG.Json, cfg.LOG.TimeFormat, false)
	grpclog.SetLoggerV2(log)

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
		// postgres.SetLogger(log),
		// postgres.SetReconnectTimeout(cfg.DB.ReconnectTimeout),
	)
	if err != nil {
		log.Errorf("connect to db: %v", err)
		return
	}

	// cli commands
	err = cli.Run(
		ctx,
		cli.NewMigrateCmd(log, d),
		cli.NewCreateEmailPassIdentityCmd(
			log,
			create_emailpass_identity.New(
				d,
				cfg.APP.BaseURL,
				cfg.APP.FromEmail,
				cfg.APP.FromName,
				cfg.APP.ConfirmLife,
			),
		),
	)
	if err != nil {
		log.Errorf("cli: %v", err)
		return
	}

	// prepare web/gRPC server handlers
	createEmailPassIdentityUseCase := create_emailpass_identity.New(
		d,
		cfg.APP.BaseURL,
		cfg.APP.FromEmail,
		cfg.APP.FromName,
		cfg.APP.ConfirmLife,
	)
	webServer := web.New(log,
		&web.Queries{
			Profile: profile.New(d),
			Whoami:  whoami.New(d, &cfg.TOKEN),
		},
		&web.Commands{
			CreateEmailPassIdentity: createEmailPassIdentityUseCase,
			AuthEmailPass:           auth_emailpass.New(d, &cfg.TOKEN),
			Confirm:                 confirm.New(d),
			RetryConfirm:            retry_confirm.New(d, createEmailPassIdentityUseCase),
		},
	)

	// run gRPC server
	grpcErr := grpc_server.Run(
		ctx,
		log,
		cfg.GRPC.Addr(), // listen incoming host:port for gRPC server
		func(s grpc.ServiceRegistrar) {
			pb.RegisterUserServiceServer(s, webServer)
		},
	)

	// run web -> gRPC gateway
	grpcGwErr := grpc_gateway.Run(
		ctx,
		log,
		cfg.HTTP.Addr(), // listen incoming host:port for rest api
		cfg.GRPC.Addr(), // connect to gRPC server host:port
		[]grpc_gateway.RegisterServiceHandlerFunc{
			pb.RegisterUserServiceHandler,
		},
	)

	log.Infof("application started (version: %s buildtime: %s)", conf.Version, conf.Buildtime)
	defer log.Infof("application finished")

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
