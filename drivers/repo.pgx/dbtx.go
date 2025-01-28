package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lucasjacques/modb"
	"github.com/lucasjacques/modb/queries"
)

type pgxdbtx interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

type dbtx struct {
	conn pgxdbtx
}

func (d *dbtx) NewParamsSet() queries.ParamsSet {
	return &queries.Numbered{}
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

func (d *dbtx) FQCN(table, column string) string {
	return `"` + table + `"."` + column + `"`
}

var _ modb.DBTX = (*dbtx)(nil)
