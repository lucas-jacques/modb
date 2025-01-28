package queries

import (
	"strings"
)

type Insert struct {
	placeholders ParamsSet
	into         string
	columns      []string
	values       [][]any
}

func (i *Insert) Values(values []any) *Insert {
	i.values = append(i.values, values)
	return i
}

func buildInsertPlaceholder(s ParamsSet, count int) string {
	builder := strings.Builder{}
	builder.WriteString("(")
	for i := 0; i < count; i++ {
		builder.WriteString(s.Next())
		if i < count-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")")

	return builder.String()
}

func (i *Insert) Build() (string, []any) {
	var values []any
	var parts []string

	parts = append(parts, "INSERT INTO", i.into)

	parts = append(parts, "("+strings.Join(i.columns, ", ")+")")
	parts = append(parts, "VALUES")

	for _, valueSet := range i.values {
		parts = append(parts, buildInsertPlaceholder(i.placeholders, len(valueSet)))
		values = append(values, valueSet...)
	}

	return strings.Join(parts, " "), values
}
