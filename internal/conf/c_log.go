package conf

import (
	"github.com/caarlos0/env"
)

type LOG struct {
	LogLevel   string `env:"FDLT_LOG_LEVEL" envDefault:"debug"`
	LogJson    bool   `env:"FDLT_LOG_JSON" envDefault:"false"`
	TimeFormat string `env:"FDLT_LOG_TIME_FORMAT" envDefault:"2006-01-02T15:04:05Z"`
}

func (c *config) LOG() *LOG {
	if c.log != nil {
		return c.log
	}

	c.log = &LOG{}

	if err := env.Parse(c.log); err != nil {
		panic("parsing LOG configuration")
	}

	return c.log
}
