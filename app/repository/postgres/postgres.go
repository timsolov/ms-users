package postgres

import (
	"context"
	"database/sql"
	"time"

	"ms-users/app/common/logger"

	"github.com/dimiro1/health"
	_ "github.com/jackc/pgx/v4/stdlib" // force include stdlib
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	_defaultConnectTimeout   = time.Second
	_defaultReconnectTimeout = time.Second
)

// DB describes
type DB struct {
	db                         *sqlx.DB
	tx                         *sqlx.Tx
	connectTimeout             time.Duration
	reconnectTimeout           time.Duration
	maxOpenConns, maxIdleConns int
	openLifeTime, idleLifeTime time.Duration
	log                        logger.Logger
}

type Option func(*DB)

func SetConnectTimeout(t time.Duration) Option {
	return func(d *DB) {
		d.connectTimeout = t
	}
}

func SetReconnectTimeout(t time.Duration) Option {
	return func(d *DB) {
		d.reconnectTimeout = t
	}
}

func SetMaxConns(maxOpen, maxIdle int) Option {
	return func(d *DB) {
		d.maxOpenConns = maxOpen
		d.maxIdleConns = maxIdle
	}
}

func SetConnsMaxLifeTime(open, idle time.Duration) Option {
	return func(d *DB) {
		d.openLifeTime = open
		d.idleLifeTime = idle
	}
}

func SetLogger(log logger.Logger) Option {
	return func(d *DB) {
		d.log = log
	}
}

func New(ctx context.Context, dsn string, opts ...Option) (*DB, error) {
	d := &DB{
		connectTimeout:   _defaultConnectTimeout,
		reconnectTimeout: _defaultReconnectTimeout,
	}
	for _, opt := range opts {
		opt(d)
	}

	var (
		db  *sqlx.DB
		err error
	)

	ctx, cancel := context.WithTimeout(ctx, d.connectTimeout)
	defer cancel()

	db, err = sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to DB: %s", dsn)
	}

	d.db = db

	if d.maxOpenConns > 0 {
		db.SetMaxOpenConns(d.maxOpenConns)
	}
	if d.maxIdleConns > 0 {
		db.SetMaxIdleConns(d.maxIdleConns)
	}
	if d.openLifeTime > 0 {
		db.SetConnMaxLifetime(d.openLifeTime)
	}
	if d.idleLifeTime > 0 {
		db.SetConnMaxIdleTime(d.idleLifeTime)
	}

	if d.reconnectTimeout > 0 {
		go d.reconnect(ctx)
	}

	return d, nil
}

func (d *DB) reconnect(ctx context.Context) {
	d.infof("postgres reconnection goroutine started")
	defer d.infof("postgres reconnection goroutine finished")

	ticker := time.NewTicker(d.reconnectTimeout)
	connected := true

	for {
		select {
		case <-ticker.C:
			err := d.db.PingContext(ctx)
			if err != nil {
				if connected {
					connected = false
					d.errorf(err, "postgres connection lost")
					continue
				}
				d.errorf(err, "postgres reconnection")
				continue
			}
			if !connected {
				connected = true
				d.infof("postgres connection established")
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (d *DB) errorf(err error, format string, args ...interface{}) {
	if d.log == nil {
		return
	}
	d.log.WithError(err).Errorf(format, args...)
}

func (d *DB) infof(format string, args ...interface{}) {
	if d.log == nil {
		return
	}
	d.log.Infof(format, args...)
}

func (d *DB) atomic(ctx context.Context, fn func(tx *DB) error) error {
	if d.tx != nil {
		return fn(d)
	}

	tx, err := d.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	copyDB := *d
	copyDB.tx = tx

	err = fn(&copyDB)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

func (d *DB) Check(ctx context.Context) health.Health {
	var res health.Health

	start := time.Now()
	_, err := d.one(ctx, "SELECT 1")
	since := time.Since(start)
	res.AddInfo("duration", since.String())
	if err != nil {
		res.AddInfo("error", err.Error())
		res.Down()
		return res
	}

	dbStats := d.db.Stats()

	stats := map[string]any{
		"idle":                 dbStats.Idle,               // The number of idle connections.
		"in_use":               dbStats.InUse,              // The number of connections currently in use.
		"max_idle_closed":      dbStats.MaxIdleClosed,      // The total number of connections closed due to SetMaxIdleConns.
		"max_idle_time_closed": dbStats.MaxIdleTimeClosed,  // The total number of connections closed due to SetConnMaxIdleTime.
		"max_life_time_closed": dbStats.MaxLifetimeClosed,  // The total number of connections closed due to SetConnMaxLifetime.
		"max_open_connections": dbStats.MaxOpenConnections, // Maximum number of open connections to the database.
		"open_connections":     dbStats.OpenConnections,    // Pool Status
		"wait_count":           dbStats.WaitCount,          // Counters
		"wait_duration":        dbStats.WaitDuration,       // The total time blocked waiting for a new connection.
	}

	res.AddInfo("stats", stats)

	res.Up()

	return res
}
