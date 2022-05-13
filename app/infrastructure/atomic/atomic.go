package atomic

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type Atomicer interface {
	Begin(ctx context.Context) (Atomicer, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func Do[Repo any](ctx context.Context, r Repo, fn func(Repo) error) error {
	atomic, ok := any(r).(Atomicer)
	if !ok {
		// return fn(db)
		return fmt.Errorf("transaction not support")
	}

	tx, err := atomic.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer atomic.Rollback(ctx)

	if err := fn(tx.(Repo)); err != nil {
		return err
	}

	if err := atomic.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}
