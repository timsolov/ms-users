package conf

import (
	"sync"

	"github.com/caarlos0/env"
)

type LOG struct {
	Level      string `env:"LOG_LEVEL" envDefault:"debug"`
	Json       bool   `env:"LOG_JSON" envDefault:"false"`
	TimeFormat string `env:"LOG_TIME_FORMAT" envDefault:"2006-01-02T15:04:05Z"`
}

var logOnce sync.Once

func (c *config) LOG() *LOG {
	logOnce.Do(func() {
		c.log = &LOG{}

		if err := env.Parse(c.log); err != nil {
			panic("parsing LOG configuration")
		}
	})

	return c.log
}
