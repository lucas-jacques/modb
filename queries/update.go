package queries

import (
	"strings"

	"github.com/lucasjacques/modb"
)

type Update struct {
	params  modb.ParamsSet
	values  []Expr
	columns []string
	table   string
	where   Expr
}

func (u *Update) Set(column string, value Expr) *Update {
	u.columns = append(u.columns, column)
	u.values = append(u.values, value)
	return u
}

func (u *Update) Where(where Expr) *Update {
	u.where = where
	return u
}

func (u *Update) Build() (string, []any) {
	builder := strings.Builder{}

	builder.WriteString("UPDATE ")
	builder.WriteString(u.table)

	builder.WriteString(" SET ")

	var values []any

	for i, col := range u.columns {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(col)
		builder.WriteString(" = ")

		part, colValues := u.values[i].Build(u.params)
		builder.WriteString(part)

		values = append(values, colValues...)
	}

	if u.where != nil {
		builder.WriteString(" WHERE ")

		part, whereValues := u.where.Build(u.params)
		builder.WriteString(part)

		values = append(values, whereValues...)
	}

	return builder.String(), values
}
