package queries

import (
	"strconv"
	"strings"

	"github.com/lucasjacques/modb"
)

type Select struct {
	params    modb.ParamsSet
	values    []any
	columns   []string
	table     string
	where     string
	limit     string
	selectFor string
	joins     []*Join
}

func NewQuery(placeholders modb.ParamsSet) *Query {
	return &Query{params: placeholders}
}

func (q *Select) Build() (string, []any) {
	var values []any
	var parts []string

	parts = append(parts, "SELECT")

	if len(q.columns) > 0 {
		parts = append(parts, strings.Join(q.columns, ", "))
	}

	parts = append(parts, "FROM", q.table)

	for _, j := range q.joins {
		sql, params := j.Build(q.params)

		parts = append(parts, sql)
		values = append(values, params...)
	}

	if q.where != "" {
		parts = append(parts, "WHERE", q.where)
		values = append(values, q.values...)
	}

	return strings.Join(parts, " "), values
}

func (q *Select) Select(table string, columns ...string) *Select {
	q.columns = columns
	q.table = table
	return q
}

func (q *Select) Limit(limit int) *Select {
	q.limit = "LIMIT " + strconv.Itoa(limit)
	return q

}

func (q *Select) Where(where Expr) *Select {
	var values []any

	q.where, values = where.Build(q.params)
	q.values = append(q.values, values...)
	return q
}

type SelectForOpt string

const (
	SkipLocked SelectForOpt = "SKIP LOCKED"
	NoWait     SelectForOpt = "NOWAIT"
)

func (s *Select) ForUpdate(opts ...SelectForOpt) *Select {
	s.selectFor = "FOR UPDATE"
	if len(opts) > 0 {
		s.selectFor += " " + string(opts[0])
	}

	return s
}

func (s *Select) ForShare(opts ...SelectForOpt) *Select {
	s.selectFor = "FOR SHARE"
	if len(opts) > 0 {
		s.selectFor += " " + string(opts[0])
	}

	return s
}

func (s *Select) ForKeyShare(opts ...SelectForOpt) *Select {
	s.selectFor = "FOR KEY SHARE"
	if len(opts) > 0 {
		s.selectFor += " " + string(opts[0])
	}
	return s
}

func (q *Select) Join(j *Join) *Select {
	q.joins = append(q.joins, j)
	return q
}
