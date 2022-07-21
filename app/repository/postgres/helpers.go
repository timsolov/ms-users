package postgres

import (
	"context"
	"database/sql"
	"ms-users/app/common/event"
	"ms-users/app/domain"

	"github.com/pkg/errors"
)

func (d *DB) publishEvents(ctx context.Context, events []event.Event) error {
	for i := 0; i < len(events); i++ {
		err := d.Publish(ctx, events[i].Subject, events[i].Payload)
		if err != nil {
			return errors.Wrapf(err, "publish events[%d]", i)
		}
	}
	return nil
}

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

func (d *DB) many(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var (
		r   *sql.Rows
		err error
	)

	if d.tx != nil {
		r, err = d.tx.QueryContext(ctx, d.db.Rebind(query), args...)
	} else {
		r, err = d.db.QueryContext(ctx, d.db.Rebind(query), args...)
	}

	return r, E(err)
}

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

// contains returns `true` when `list` includes `v`.
func contains(list []string, v string) bool {
	for _, s := range list {
		if v == s {
			return true
		}
	}
	return false
}

type containsAlgo uint8

const (
	atLeastOne containsAlgo = iota
	exactlyAll
)

// mcontains returns `true` when `list` includes all `values`.
// `algo` means algorithm for matching.
// `atLeastOne` - one value from values enough for positive;
// `exactlyAll` - exectly all values should be in the list.
func mcontains(list, values []string, algo containsAlgo) bool {
	for _, v := range values {
		present := contains(list, v)
		if algo == exactlyAll && !present {
			return false
		}
		if algo == atLeastOne && present {
			return true
		}
	}

	return algo == exactlyAll
}
