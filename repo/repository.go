package repo

import (
	"context"
	"fmt"

	"github.com/lucasjacques/modb"
	"github.com/lucasjacques/modb/model"
	"github.com/lucasjacques/modb/queries"
)

func New[M any, PK comparable, C model.Schema[M]](db modb.DBTX, model model.TypedModelCols[M, PK, C]) *ModelRepository[M, PK, C] {
	return &ModelRepository[M, PK, C]{
		db:    db,
		model: model,
	}
}

type ModelRepository[M any, PK comparable, C model.Schema[M]] struct {
	db    modb.DBTX
	model model.TypedModelCols[M, PK, C]
}

func (m *ModelRepository[M, PK, S]) Insert(ctx context.Context, model *M) error {
	table := m.model.GetTable()
	cols := m.model.GetColumns()

	columns := make([]string, 0, len(cols))
	values := make([]any, 0, len(cols))

	for _, col := range cols {
		if col.ShouldOmit(m, modb.OpInsert) {
			continue
		}

		columns = append(columns, col.GetName())
		value, err := col.ValueFromModel(model)
		if err != nil {
			return err
		}

		values = append(values, value)
	}

	str, values := queries.NewQuery(m.db.NewParamsSet()).Insert(table, columns).Values(values).Build()

	_, err := m.db.Exec(ctx, str, values...)
	if err != nil {
		return err
	}

	return nil
}

func (m *ModelRepository[M, PK, C]) Update(ctx context.Context, model *M) error {
	table := m.model.GetTable()

	query := queries.NewQuery(m.db.NewParamsSet()).Update(table)

	for _, col := range m.model.GetColumns() {
		if col.ShouldOmit(model, modb.OpUpdate) {
			continue
		}

		value, err := col.ValueFromModel(model)
		if err != nil {
			return err
		}

		query.Set(col.GetName(), queries.Value(value))
	}

	pk, err := m.model.PrimaryKey().ValueFromModel(model)
	if err != nil {
		return err
	}

	sql, values := query.Where(queries.EQ(m.model.PrimaryKey(), queries.Value(pk))).Build()
	fmt.Println(sql)
	_, err = m.db.Exec(ctx, sql, values...)
	if err != nil {
		return err
	}

	return nil
}
