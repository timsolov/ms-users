package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ms-users/app/common/event"
	"ms-users/app/common/jsonschema"
	"ms-users/app/common/logger"
	"ms-users/app/conf"
	"ms-users/app/delivery/cli"
	"ms-users/app/delivery/grpc_gateway"
	"ms-users/app/delivery/grpc_listener"
	"ms-users/app/delivery/grpc_server"
	"ms-users/app/delivery/grpc_server/pb"
	"ms-users/app/repository/grpc_local"
	"ms-users/app/repository/postgres"
	"ms-users/app/usecase/auth_emailpass"
	"ms-users/app/usecase/confirm"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"
	"ms-users/app/usecase/reset_password_confirm"
	"ms-users/app/usecase/reset_password_init"
	"ms-users/app/usecase/retry_confirm"
	"ms-users/app/usecase/update_profile"
	"ms-users/app/usecase/whoami"

	"github.com/dimiro1/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	grpc_gateway.Status = health.NotReady

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

	// jsonschema for profile
	jsonSchema := jsonschema.New([]string{"properties.avatar_url"})
	err = jsonSchema.Load(ctx, cfg.APP.ProfileJSONSchemaPath)
	if err != nil {
		log.Errorf("load jsonschema: %v", err)
		return
	}
	jsonSchemaName := strings.TrimSuffix(filepath.Base(cfg.APP.ProfileJSONSchemaPath), filepath.Ext(cfg.APP.ProfileJSONSchemaPath))

	// cli commands
	err = cli.Run(
		ctx,
		cli.NewMigrateCmd(log, d),
		cli.NewCreateEmailPassIdentityCmd(
			log,
			create_emailpass_identity.New(
				d,
				cfg.APP.APIBaseURL,
				cfg.APP.FromEmail,
				cfg.APP.FromName,
				cfg.APP.ConfirmLife,
				jsonSchema,
				jsonSchemaName,
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
		cfg.APP.APIBaseURL,
		cfg.APP.FromEmail,
		cfg.APP.FromName,
		cfg.APP.ConfirmLife,
		jsonSchema,
		jsonSchemaName,
	)
	grpcServer := grpc_server.New(log,
		&grpc_server.Queries{
			Profile: profile.New(d),
			Whoami:  whoami.New(d, &cfg.TOKEN),
		},
		&grpc_server.Commands{
			CreateEmailPassIdentity: createEmailPassIdentityUseCase,
			AuthEmailPass:           auth_emailpass.New(d, &cfg.TOKEN),
			Confirm:                 confirm.New(d),
			RetryConfirm:            retry_confirm.New(d, createEmailPassIdentityUseCase),
			ResetPasswordInit: reset_password_init.New(
				d,
				cfg.APP.WebBaseURL,
				cfg.APP.FromEmail,
				cfg.APP.FromName,
				cfg.APP.ConfirmLife,
			),
			ResetPasswordConfirm: reset_password_confirm.New(d),
			UpdateProfile:        update_profile.New(d, jsonSchema, createEmailPassIdentityUseCase),
		},
	)

	// run gRPC server
	grpcErr := grpc_listener.Run(
		ctx,
		log,
		cfg.GRPC.Addr(), // listen incoming host:port for gRPC server
		func(s grpc.ServiceRegistrar) {
			pb.RegisterUserServiceServer(s, grpcServer)
			grpc_health_v1.RegisterHealthServer(s, grpcServer)
		},
	)

	// gRPC health client
	grpcHealthClient, err := grpc_local.NewHealthClient(ctx, cfg.GRPC.Addr())
	if err != nil {
		log.Errorf("health client: %v", err)
		return
	}

	// healthCheck
	healthCheck := health.NewHandler()
	healthCheck.AddChecker("ms-email", event.HealthChecker(cfg.EMAIL.Addr()))
	healthCheck.AddChecker("postgres", d)
	healthCheck.AddChecker("grpc_server", grpc_local.HealthChecker(grpcHealthClient))
	healthCheck.AddInfo("app", map[string]any{
		"version":   conf.Version,
		"buildtime": conf.Buildtime,
	})

	// run web -> gRPC gateway
	grpcGwErr := grpc_gateway.Run(
		ctx,
		log,
		cfg.HTTP.Addr(), // listen incoming host:port for rest api
		cfg.GRPC.Addr(), // connect to gRPC server host:port
		cfg.HTTP.Timeout,
		healthCheck,
		[]grpc_gateway.RegisterServiceHandlerFunc{
			pb.RegisterUserServiceHandler,
		},
	)

	log.Infof("application started (version: %s buildtime: %s)", conf.Version, conf.Buildtime)
	defer log.Infof("application finished")

	grpc_gateway.Status = health.Up

	select {
	case <-ctx.Done():
	case err := <-grpcErr:
		log.Errorf("gRPC server error: %s", err)
	case err := <-grpcGwErr:
		log.Errorf("gRPC gateway error: %s", err)
	}

	grpc_gateway.Status = health.Down

	cancel()

	time.Sleep(1 * time.Second)
}
