package conf

import (
	"sync"

	"github.com/caarlos0/env"
)

// APP describes
type APP struct {
}

var appOnce sync.Once

func (c *config) APP() *APP {
	appOnce.Do(func() {
		c.app = &APP{}

		if err := env.Parse(c.app); err != nil {
			panic("parsing APP configuration")
		}
	})

	return c.app
}
