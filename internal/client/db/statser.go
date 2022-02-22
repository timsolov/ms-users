package db

import (
	"database/sql"
	"fmt"
)

type statser interface {
	// Stats returns database statistics.
	Stats() sql.DBStats
}

func Stats(d DB) (sql.DBStats, error) {
	statser, ok := d.(statser)
	if !ok {
		return sql.DBStats{}, fmt.Errorf("Stats method not supported")
	}
	return statser.Stats(), nil
}
