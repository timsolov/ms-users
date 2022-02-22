package conf

import "sync"

type Config interface {
	APP() *APP
	DB() *DB
	HTTP() *HTTP
	LOG() *LOG
}

func New() Config {
	c := &config{}
	c.APP()
	c.DB()
	c.HTTP()
	c.LOG()
	return c
}

type config struct {
	mu   sync.Mutex
	app  *APP
	db   *DB
	http *HTTP
	log  *LOG
}
