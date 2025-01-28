package repo

import (
	"context"

	"github.com/lucasjacques/modb"
	"github.com/lucasjacques/modb/queries"
)

type dbtx struct {
	c stddbtx
}

var _ modb.DBTX = (*dbtx)(nil)

func (d *dbtx) Exec(ctx context.Context, query string, args ...any) (modb.CommandTag, error) {
	return d.c.ExecContext(ctx, query, args...)
}

func (d *dbtx) Query(ctx context.Context, query string, args ...any) (modb.Rows, error) {
	return d.c.QueryContext(ctx, query, args...)
}

func (d *dbtx) NewParamsSet() queries.ParamsSet {
	return &queries.QuestionMark{}
}

func (d *dbtx) FQCN(table, column string) string {
	return `"` + table + `"."` + column + `"`
}

func wrapDBTX(conn stddbtx) modb.DBTX {
	return &dbtx{
		c: conn,
	}
}
