package grpc_local

import (
	"context"
	"time"

	"github.com/dimiro1/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewHealthClient(ctx context.Context, target string) (grpc_health_v1.HealthClient, error) {
	conn, err := grpc.DialContext(ctx, target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return grpc_health_v1.NewHealthClient(conn), nil
}

func HealthChecker(grpcHealthClient grpc_health_v1.HealthClient) health.CheckerFunc {
	return func(ctx context.Context) health.Health {
		var res health.Health
		res.Down()

		start := time.Now()

		response, err := grpcHealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			res.AddInfo("error", err.Error())
			return res
		}

		duration := time.Since(start)
		res.AddInfo("duration", duration.String())

		if response.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			res.Up()
		}

		return res
	}
}
