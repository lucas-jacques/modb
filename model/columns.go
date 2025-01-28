package model

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/lucasjacques/modb/queries"
)

type column[M any, V any] struct {
	table        string
	name         string
	innerModel   TypedModel[M]
	ref          func(*M) *V
	omitOnInsert bool
	omitOnUpdate bool
}

// New implements TypedCol.
func (c *column[M, V]) New() *V {
	return new(V)
}

// ShouldOmit implements TypedCol.
func (c *column[M, V]) ShouldOmit(m any, op Operation) bool {
	switch op {
	case OpInsert:
		return c.omitOnInsert
	case OpUpdate:
		return c.omitOnUpdate
	default:
		return false
	}
}

var _ TypedCol[any, any] = (*column[any, any])(nil)

func (c *column[M, V]) setTable(table string) {
	c.table = table
}

func (c *column[M, V]) SetInnerModel(innerModel TypedModel[M]) {
	c.innerModel = innerModel
}

func (c *column[M, V]) GetInnerModel() TypedModel[M] {
	return c.innerModel
}

func (c *column[M, V]) ForDest(m any) (any, error) {
	model, ok := m.(*M)
	if !ok {
		return nil, ErrWrongModelType
	}

	return c.ref(model), nil
}

// Implement modb.Column[M, V]

func (c *column[M, V]) GetTable() string {
	return c.table
}

func (c *column[M, V]) GetName() string {
	return c.name
}

func (c *column[M, V]) GetOmitOnInsert() bool {
	return c.omitOnInsert
}

func (c *column[M, V]) GetOmitOnUpdate() bool {
	return c.omitOnUpdate
}

func (c *column[M, V]) NewDest() any {
	return new(V)
}

func (c *column[M, V]) SetValueOnModel(m any, v any) error {
	ptr, ok := v.(*V)
	if !ok {
		return ErrWrongValueType
	}

	model, ok := m.(*M)
	if !ok {
		return ErrWrongModelType
	}

	ref := c.ref(model)
	if ref == nil {
		return errors.New("invalid model")
	}

	if ptr == nil {
		return nil
	}

	*ref = *ptr
	return nil
}

var ErrWrongValueType = fmt.Errorf("invalid value type")

var ErrWrongModelType = fmt.Errorf("invalid model type")

func (c *column[M, V]) ValueFromModel(m any) (any, error) {
	model, ok := m.(*M)
	if !ok {
		fmt.Println("model", reflect.TypeOf(m).String())
		return nil, ErrWrongModelType
	}

	return *c.ref(model), nil
}

func (c *column[M, V]) ValueFromModelTyped(m *M) (any, error) {
	return *c.ref(m), nil
}

func (c *column[M, V]) GetValue(m *M) (V, error) {
	return *c.ref(m), nil
}

// Options

func (c *column[M, V]) OmitOnInsert() *column[M, V] {
	c.omitOnInsert = true
	return c
}

func (c *column[M, V]) OmitOnUpdate() *column[M, V] {
	c.omitOnUpdate = true
	return c
}

func (c *column[M, V]) FQCN() string {
	return `"` + c.table + `".` + `"` + c.name + `"`
}

func (c *column[M, V]) Build(queries.ParamsSet) (string, []any) {
	return c.FQCN(), nil
}

func Col[M, V any](name string, ref func(*M) *V) *column[M, V] {
	return &column[M, V]{name: name, ref: ref}
}

func PrimaryKey[M any, V any](name string, ref func(*M) *V) *column[M, V] {
	return &column[M, V]{
		name:         name,
		ref:          ref,
		omitOnUpdate: true,
	}
}

type autoincrement interface {
	~int | int64 | ~uint | ~uint64 | ~int32 | ~uint32
}

func AutoIncrement[M any, V autoincrement](name string, ref func(*M) *V) *column[M, V] {
	return &column[M, V]{
		name:         name,
		ref:          ref,
		omitOnInsert: true,
		omitOnUpdate: true,
	}
}
