package modb

import (
	"context"
	"database/sql"
	"errors"
)

type Scannable interface {
	Scan(dest ...any) error
}

type DBTX interface {
	Exec(ctx context.Context, query string, args ...any) (CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	NewParamsSet() ParamsSet
}

type TX interface {
	DBTX
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxOptions = sql.TxOptions

type Database interface {
	DBTX
	Begin(ctx context.Context) (TX, error)
	BeginTx(ctx context.Context, opts *TxOptions) (TX, error)
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

var ErrNoRows = errors.New("no rows in result set")
