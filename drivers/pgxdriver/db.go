package pgxdriver

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/lucasjacques/modb"
)

type PGXConn interface {
	pgxdbtx
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}

type DB struct {
	conn PGXConn
	dbtx
}

func NewMODB(conn PGXConn) modb.Database {
	return &DB{
		conn: conn,
		dbtx: dbtx{conn: conn},
	}
}

type tx struct {
	dbtx
	tx pgx.Tx
}

// Commit implements modb.TX.
func (t *tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

// Rollback implements modb.TX.
func (t *tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

var _ modb.TX = (*tx)(nil)

// Begin implements modb.Database.
func (d *DB) Begin(ctx context.Context) (modb.TX, error) {
	t, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &tx{
		dbtx: dbtx{
			conn: t,
		},
		tx: t,
	}, nil
}

func getIsolationLevel(iso sql.IsolationLevel) pgx.TxIsoLevel {
	switch iso {
	case sql.LevelReadUncommitted:
		return pgx.ReadUncommitted
	case sql.LevelReadCommitted:
		return pgx.ReadCommitted
	case sql.LevelRepeatableRead:
		return pgx.RepeatableRead
	case sql.LevelSerializable:
		return pgx.Serializable
	}
	return pgx.ReadCommitted
}

// BeginTx implements modb.Database.
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (modb.TX, error) {
	var accessMode pgx.TxAccessMode
	if opts.ReadOnly {
		accessMode = pgx.ReadOnly
	} else {
		accessMode = pgx.ReadWrite
	}

	t, err := d.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   getIsolationLevel(opts.Isolation),
		AccessMode: accessMode,
	})
	if err != nil {
		return nil, err
	}

	return &tx{
		dbtx: dbtx{
			conn: t,
		},
		tx: t,
	}, nil

}

var _ modb.Database = &DB{}
