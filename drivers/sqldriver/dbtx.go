package sqldriver

import (
	"context"

	"github.com/lucasjacques/modb"
)

type dbtx struct {
	fqcn FQCNMode
	c    stddbtx
}

var _ modb.DBTX = (*dbtx)(nil)

// QueryRow implements modb.DBTX.
func (d *dbtx) QueryRow(ctx context.Context, query string, args ...any) modb.Row {
	return d.c.QueryRowContext(ctx, query, args...)
}

func (d *dbtx) Exec(ctx context.Context, query string, args ...any) (modb.CommandTag, error) {
	return d.c.ExecContext(ctx, query, args...)
}

func (d *dbtx) Query(ctx context.Context, query string, args ...any) (modb.Rows, error) {
	return d.c.QueryContext(ctx, query, args...)
}

func (d *dbtx) NewParamsSet() modb.ParamsSet {
	return &modb.QuestionMark{}
}

func (d *dbtx) FQCN(table, column string) string {
	switch d.fqcn {
	case FQCNDoubleQuotes:
		return `"` + table + `"."` + column + `"`
	case FQCNBackticks:
		return "`" + table + "`." + "`" + column + "`"
	}

	return column
}
