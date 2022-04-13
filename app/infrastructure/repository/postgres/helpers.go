package postgres

import (
	"context"
	"database/sql"

	"github.com/timsolov/ms-users/app/domain/entity"
)

func (d *DB) get(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	return d.db.GetContext(ctx, dst, d.db.Rebind(query), args...)
}

func (d *DB) execr(ctx context.Context, rows int64, query string, args ...interface{}) error {
	r, err := d.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	if i, err := r.RowsAffected(); err != nil {
		return err
	} else if i != rows {
		return entity.ErrMismatch
	}

	return nil
}

func E(err error) error {
	if err == sql.ErrNoRows {
		return entity.ErrNotFound
	}
	return err
}
