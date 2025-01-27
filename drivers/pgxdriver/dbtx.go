package pgxdriver

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lucasjacques/modb"
)

type row struct {
	err  error
	rows pgx.Rows
}

func (r *row) Err() error {
	return r.err
}

func (r *row) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	if !r.rows.Next() {
		r.rows.Close()
		return modb.ErrNoRows
	}

	return r.rows.Scan(dest...)
}

type pgxdbtx interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type dbtx struct {
	conn pgxdbtx
}

func (d *dbtx) NewParamsSet() modb.ParamsSet {
	return &modb.Numbered{}
}

func (c *cmdTag) RowsAffected() (int64, error) {
	return c.tag.RowsAffected(), nil
}

type cmdTag struct {
	tag pgconn.CommandTag
}

func (d *dbtx) Exec(ctx context.Context, query string, args ...any) (modb.CommandTag, error) {
	tag, err := d.conn.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &cmdTag{tag: tag}, nil
}

func (d *dbtx) Query(ctx context.Context, query string, args ...any) (modb.Rows, error) {
	r, err := d.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &rows{rows: r}, nil
}

func (d *dbtx) QueryRow(ctx context.Context, query string, args ...any) modb.Row {
	rows, err := d.conn.Query(ctx, query, args...)
	if err != nil {
		return &row{err: err}
	}

	return &row{rows: rows}
}

func (d *dbtx) FQCN(table, column string) string {
	return `"` + table + `"."` + column + `"`
}

var _ modb.DBTX = (*dbtx)(nil)
