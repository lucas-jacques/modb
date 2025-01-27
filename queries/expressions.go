package queries

import (
	"github.com/lucasjacques/modb"
)

type Expr interface {
	Build(modb.ParamsSet) (string, []any)
}

type value struct {
	val any
}

func (v *value) Build(p modb.ParamsSet) (string, []any) {
	return p.Next(), []any{v.val}
}

func Value(val any) Expr {
	return &value{val}
}
