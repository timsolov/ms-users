package conf

import (
	"sync"

	"github.com/caarlos0/env"
)

type GRPC struct {
	Host string `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	Port string `env:"GRPC_PORT" envDefault:"10000"`
}

func (grpc *GRPC) Addr() string {
	return grpc.Host + ":" + grpc.Port
}

var gRPCOnce sync.Once

func (c *config) GRPC() *GRPC {
	gRPCOnce.Do(func() {
		c.grpc = &GRPC{}

		if err := env.Parse(c.grpc); err != nil {
			panic("parsing GRPC configuration")
		}
	})

	return c.grpc
}
