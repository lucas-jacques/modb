package queries

import (
	"strings"

	"github.com/lucasjacques/modb"
)

type Join struct {
	join  string // JOIN, LEFT JOIN, RIGHT JOIN, etc.
	table string
	on    Expr
}

func (j *Join) Build(p modb.ParamsSet) (string, []any) {
	var parts []string
	var params []any

	parts = append(parts, j.join, j.table)

	if j.on != nil {
		sql, values := j.on.Build(p)
		parts = append(parts, "ON", sql)
		params = append(params, values...)
	}

	return strings.Join(parts, " "), params
}

func LeftJoin(table string) *Join {
	return &Join{join: "LEFT JOIN", table: table}
}

func RightJoin(table string) *Join {
	return &Join{join: "RIGHT JOIN", table: table}
}

func InnerJoin(table string) *Join {
	return &Join{join: "INNER JOIN", table: table}
}

type raw string

func (r raw) Build(_ modb.ParamsSet) (string, []any) {
	return string(r), nil
}

func Raw(expr string) Expr {
	return raw(expr)
}

func (j *Join) On(expr Expr) *Join {
	j.on = expr
	return j
}
