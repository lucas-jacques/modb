package sqldriver

import (
	"context"
	"database/sql"

	"github.com/lucasjacques/modb"
)

type db struct {
	dbtx
	db *sql.DB
}

type FQCNMode int

const (
	FQCNDoubleQuotes FQCNMode = iota
	FQCNBackticks
)

func NewMODB(database *sql.DB, fqcn FQCNMode) modb.Database {
	return &db{
		dbtx: dbtx{c: database, fqcn: fqcn},
		db:   database,
	}
}

type tx struct {
	dbtx
	tx *sql.Tx
}

// Commit implements modb.TX.
func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

// Rollback implements modb.TX.
func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}

var _ modb.TX = (*tx)(nil)

func (db *db) Begin(ctx context.Context) (modb.TX, error) {
	t, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &tx{
		tx:   t,
		dbtx: dbtx{c: t, fqcn: db.fqcn},
	}, nil
}

type stddbtx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// BeginTx implements modb.Database.
func (d *db) BeginTx(ctx context.Context, opts *modb.TxOptions) (modb.TX, error) {
	t, err := d.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &tx{
		dbtx: dbtx{c: t, fqcn: d.fqcn},
		tx:   t,
	}, nil
}
