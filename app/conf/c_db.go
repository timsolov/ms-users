package conf

import (
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env"
)

// DB describes database config
type DB struct {
	Name     string `env:"DB_NAME,required"`
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	SSL      string `env:"DB_SSL,required"`
	TimeZone string `env:"DB_TIMEZONE,required"`

	OpenLimit int           `env:"DB_OPEN_LIMIT" envDefault:"5"`
	IdleLimit int           `env:"DB_IDLE_LIMIT" envDefault:"5"`
	ConnLife  time.Duration `env:"DB_CONN_LIFE" envDefault:"5m"`
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
