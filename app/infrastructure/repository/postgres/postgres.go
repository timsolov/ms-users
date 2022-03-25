package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/domain/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func New(ctx context.Context, dsn string, maxConns, maxIdle int, connLifeTime time.Duration) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	d, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connections Pool settings:
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	d.SetMaxOpenConns(maxConns)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	d.SetMaxIdleConns(maxIdle)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	d.SetConnMaxLifetime(connLifeTime)

	if err := d.PingContext(ctx); err != nil {
		return nil, err
	}

	return &DB{db: db}, err
}

func (d *DB) SqlDB() (*sql.DB, error) {
	return d.db.DB()
}

// Stats returns database statistics.
func (d *DB) Stats() (stats sql.DBStats) {
	db, err := d.SqlDB()
	if err != nil {
		return
	}
	return db.Stats()
}

func (d *DB) GormDB() *gorm.DB {
	return d.db
}

func (d *DB) Atomic(ctx context.Context, fn func(r repository.Repository) error) error {
	db := d.db.WithContext(ctx)

	tx := db.Begin()
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "begin tx")
	}
	defer tx.Rollback()

	if err := fn(&DB{tx}); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

// E helper function to replace specific driver related NotFound error to generic between all db drivers.
// Other errors will be returned without replacing. (gorm.ErrRecordNotFound -> entity.ErrNotFound, other error -> other error)
func E(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.ErrNotFound
	}
	if IsUniqueViolationErr(err) {
		return entity.ErrNotUnique
	}
	return err
}
