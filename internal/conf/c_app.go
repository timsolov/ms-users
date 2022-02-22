package conf

import "github.com/caarlos0/env"

// APP describes
type APP struct {
}

func (c *config) APP() *APP {
	if c.app != nil {
		return c.app
	}

	c.app = &APP{}

	if err := env.Parse(c.app); err != nil {
		panic("parsing APP configuration")
	}

	return c.app
}
