package conf

import (
	"sync"

	"github.com/caarlos0/env"
)

type HTTP struct {
	Host string   `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port string   `env:"HTTP_PORT" envDefault:"11000"`
	CORS []string `env:"HTTP_CORS" envSeparator:";" envDefault:"*"`
}

func (http *HTTP) Addr() string {
	return http.Host + ":" + http.Port
}

var httpOnce sync.Once

func (c *config) HTTP() *HTTP {
	httpOnce.Do(func() {
		c.http = &HTTP{}

		if err := env.Parse(c.http); err != nil {
			panic("parsing HTTP configuration")
		}
	})

	return c.http
}
