package postgres

import (
	"context"
	"database/sql"
	"ms-users/app/domain"
)

func (d *DB) execr(ctx context.Context, rows int64, query string, args ...interface{}) error { //nolint: unparam
	var (
		r   sql.Result
		err error
	)

	if d.tx != nil {
		r, err = d.tx.ExecContext(ctx, d.db.Rebind(query), args...)
	} else {
		r, err = d.db.ExecContext(ctx, d.db.Rebind(query), args...)
	}
	if err != nil {
		return E(err)
	}

	if rows > 0 {
		i, err := r.RowsAffected()
		if err != nil {
			return E(err)
		}
		if i != rows {
			return domain.ErrMismatch
		}
	}

	return nil
}

// func (d *DB) many(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
// 	var (
// 		r   *sql.Rows
// 		err error
// 	)

// 	if d.tx != nil {
// 		r, err = d.tx.QueryContext(ctx, d.db.Rebind(query), args...)
// 	} else {
// 		r, err = d.db.QueryContext(ctx, d.db.Rebind(query), args...)
// 	}

// 	return r, E(err)
// }

// one makes sql query and returns one row from db.
// It do E(err).
func (d *DB) one(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	var (
		r *sql.Row
	)

	if d.tx != nil {
		r = d.tx.QueryRowContext(ctx, d.db.Rebind(query), args...)
	} else {
		r = d.db.QueryRowContext(ctx, d.db.Rebind(query), args...)
	}

	return r, E(r.Err())
}
