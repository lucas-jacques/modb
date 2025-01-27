package modb

type Expr interface {
	Build(ParamsSet) (string, []any)
}

type Operation uint8

const (
	OpInsert Operation = iota
	OpUpdate
	OpSelect
	OpDelete
)
