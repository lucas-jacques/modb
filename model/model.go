package model

import (
	"github.com/lucasjacques/modb"
)

type Schema[M any] interface {
	Cols() []ModelCol[M]
}

type ModelCols[M any, PK comparable] interface {
	Schema[M]
	PrimaryKey() TypedCol[M, PK]
}

type Column interface {
	GetName() string
	GetTable() string
	NewDest() any
	ForDest(m any) (any, error)
	ShouldOmit(m any, op modb.Operation) bool
	SetValueOnModel(m any, v any) error
	ValueFromModel(m any) (any, error)
	FQCN() string
	Build(modb.ParamsSet) (string, []any)
	setTable(string)
}

type ModelCol[M any] interface {
	Column
	ValueFromModelTyped(m *M) (any, error)
	GetInnerModel() TypedModel[M]
	SetInnerModel(m TypedModel[M])
}

type TypedCol[M any, V any] interface {
	ModelCol[M]
	New() *V
}

type ValueTypedCol[V any] interface {
	Column
	New() *V
}

type Model interface {
	GetTable() string
	GetPrimaryKey() Column
	GetColumns() []Column
	NewDests() []any
	FromDests(dests []any) (any, error)
}

type TypedModel[M any] interface {
	Model
	New() *M
	FromDestsTyped(dests []any) (*M, error)
}

type TypedModelCols[M any, PK comparable, C Schema[M]] interface {
	TypedModel[M]
	PrimaryKey() TypedCol[M, PK]
	Cols() C
}

func New[M any, C ModelCols[M, PK], PK comparable](table string, columns C) TypedModelCols[M, PK, C] {

	m := &model[M, PK, C]{
		table:      table,
		primaryKey: columns.PrimaryKey(),
		schema:     columns,
	}

	for _, col := range columns.Cols() {
		col.setTable(table)
		col.SetInnerModel(m)
	}

	return m
}

type model[M any, PK comparable, C Schema[M]] struct {
	table      string
	schema     C
	primaryKey TypedCol[M, PK]
}

var _ Model = (*model[any, int, Schema[any]])(nil)

func (m *model[M, PK, C]) GetTable() string {
	return m.table
}

func (m *model[M, PK, C]) Cols() C {
	return m.schema
}

func (m *model[M, PK, C]) GetColumns() []Column {
	cols := make([]Column, 0, len(m.schema.Cols()))
	for _, col := range m.schema.Cols() {
		cols = append(cols, col)
	}
	return cols
}

func (m *model[M, PK, C]) PrimaryKey() TypedCol[M, PK] {
	return m.primaryKey
}

func (m *model[M, PK, C]) GetPrimaryKey() Column {
	return m.primaryKey
}

func (m *model[M, PK, C]) New() *M {
	return new(M)
}

func (m *model[M, PK, C]) NewDests() []any {
	var dests []any
	cols := m.GetColumns()
	for _, col := range cols {
		dests = append(dests, col.NewDest())
	}
	return dests
}

func (m *model[M, PK, C]) FromDests(dests []any) (any, error) {
	model, err := m.FromDestsTyped(dests)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (m *model[M, PK, C]) FromDestsTyped(dests []any) (*M, error) {
	model := m.New()
	for i, col := range m.GetColumns() {
		if err := col.SetValueOnModel(model, dests[i]); err != nil {
			return nil, err
		}
	}
	return model, nil
}
