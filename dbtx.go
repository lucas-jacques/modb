package modb

import (
	"context"

	"github.com/lucasjacques/modb/queries"
)

type Scannable interface {
	Scan(dest ...any) error
}

type DBTX interface {
	Exec(ctx context.Context, query string, args ...any) (CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	NewParamsSet() queries.ParamsSet
}

type CommandTag interface {
	RowsAffected() (int64, error)
}

type Rows interface {
	Scannable
	Err() error
	Close() error
	Next() bool
}

type Row interface {
	Scannable
	Err() error
}
