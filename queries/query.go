package queries

import "github.com/lucasjacques/modb"

type Query struct {
	params modb.ParamsSet
}

func (q *Query) Select(table string, columns ...string) *Select {
	return &Select{
		params:  q.params,
		table:   table,
		columns: columns,
	}
}

func (q *Query) Insert(into string, columns []string) *Insert {
	return &Insert{
		placeholders: q.params,
		into:         into,
		columns:      columns,
	}
}

func (q *Query) Update(table string) *Update {
	return &Update{
		params: q.params,
		table:  table,
	}
}
