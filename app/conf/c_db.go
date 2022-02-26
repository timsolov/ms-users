package conf

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env"
)

// DB describes database config
type DB struct {
	Name     string `env:"FDLT_DATABASE_NAME,required"`
	Host     string `env:"FDLT_DATABASE_HOST,required"`
	Port     int    `env:"FDLT_DATABASE_PORT,required"`
	User     string `env:"FDLT_DATABASE_USER,required"`
	Password string `env:"FDLT_DATABASE_PASSWORD,required"`
	SSL      string `env:"FDLT_DATABASE_SSL,required"`
	TimeZone string `env:"FDLT_DATABASE_TIMEZONE,required"`
}

func (d *DB) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s timezone=%s", d.Host, d.Port, d.User, d.Password, d.Name, d.SSL, d.TimeZone)
}

var dbOnce sync.Once

func (c *config) DB() *DB {
	dbOnce.Do(func() {
		c.db = &DB{}

		if err := env.Parse(c.db); err != nil {
			panic("parsing DB configuration")
		}
	})

	return c.db
}
