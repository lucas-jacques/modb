package queries

type Expr interface {
	Build(ParamsSet) (string, []any)
}

type value struct {
	val any
}

func (v *value) Build(p ParamsSet) (string, []any) {
	return p.Next(), []any{v.val}
}

func Value(val any) Expr {
	return &value{val}
}
