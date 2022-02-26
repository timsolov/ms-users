package conf

type Config interface {
	APP() *APP
	DB() *DB
	GRPC() *GRPC
	HTTP() *HTTP
	LOG() *LOG
}

func New() Config {
	c := &config{}
	c.APP()
	c.DB()
	c.GRPC()
	c.HTTP()
	c.LOG()
	return c
}

type config struct {
	app  *APP
	db   *DB
	grpc *GRPC
	http *HTTP
	log  *LOG
}
