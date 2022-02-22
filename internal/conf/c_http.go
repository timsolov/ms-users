package conf

import (
	"github.com/caarlos0/env"
)

type HTTP struct {
	Host         string   `env:"FDLT_API_HOST" envDefault:"localhost"`
	PublicPort   string   `env:"FDLT_API_PORT" envDefault:"18080"`
	InternalPort string   `env:"FDLT_API_INTERNAL_PORT" envDefault:"8081"`
	CORS         []string `env:"FDLT_API_CORS" envSeparator:";" envDefault:"*"`
}

func (http *HTTP) PublicAddr() string {
	return http.Host + ":" + http.PublicPort
}

func (c *config) HTTP() *HTTP {
	if c.http != nil {
		return c.http
	}

	c.http = &HTTP{}

	if err := env.Parse(c.http); err != nil {
		panic("parsing HTTP configuration")
	}

	return c.http
}
