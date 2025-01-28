package modb

import (
	"context"

	"github.com/lucasjacques/modb/model"
	"github.com/lucasjacques/modb/queries"
)

func NewRepository[M any, PK comparable, C any](db DBTX, model model.TypedModelCols[M, PK, C]) *ModelRepository[M, PK, C] {
	return &ModelRepository[M, PK, C]{
		db:    db,
		model: model,
	}
}

type ModelRepository[M any, PK comparable, C any] struct {
	db    DBTX
	model model.TypedModelCols[M, PK, C]
}

func (r *ModelRepository[M, PK, S]) Insert(ctx context.Context, m *M) error {
	table := r.model.GetTable()
	cols := r.model.GetColumns()

	columns := make([]string, 0, len(cols))
	values := make([]any, 0, len(cols))

	for _, col := range cols {
		if col.ShouldOmit(r, model.OpInsert) {
			continue
		}

		columns = append(columns, col.GetName())
		value, err := col.ValueFromModel(m)
		if err != nil {
			return err
		}

		values = append(values, value)
	}

	str, values := queries.NewQuery(r.db.NewParamsSet()).Insert(table, columns).Values(values).Build()
	_, err := r.db.Exec(ctx, str, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ModelRepository[M, PK, C]) Update(ctx context.Context, m *M) error {
	table := r.model.GetTable()

	query := queries.NewQuery(r.db.NewParamsSet()).Update(table)

	for _, col := range r.model.GetColumns() {
		if col.ShouldOmit(m, model.OpUpdate) {
			continue
		}

		value, err := col.ValueFromModel(m)
		if err != nil {
			return err
		}

		query.Set(col.GetName(), queries.Value(value))
	}

	pk, err := r.model.PrimaryKey().ValueFromModel(m)
	if err != nil {
		return err
	}

	sql, values := query.Where(queries.EQ(r.model.PrimaryKey(), queries.Value(pk))).Build()
	_, err = r.db.Exec(ctx, sql, values...)
	if err != nil {
		return err
	}

	return nil
}
