package model

type Relation interface {
	ForeignDef() Model
	LocalKey() Column
}

type oneToOne[M, R any] struct {
	self       TypedModel[M]
	rel        TypedModel[R]
	foreignKey Column
	localKey   Column
	ref        func(*M) *R
}

func (oo *oneToOne[M, R]) set(m *M, r *R) {
	ref := oo.ref(m)
	*ref = *r
}

func (oo *oneToOne[M, R]) Set(m any, r any) {
	model, ok := m.(*M)
	if !ok {
		panic("invalid model type")
	}

	rel, ok := r.(*R)
	if !ok {
		panic("invalid model type")
	}

	oo.set(model, rel)
}
func (oo *oneToOne[M, R]) WithLocalKey(localKey ModelCol[M]) *oneToOne[M, R] {
	oo.localKey = localKey
	return oo
}

func (oo *oneToOne[M, R]) ForeignKey() Column {
	return oo.foreignKey
}

func (oo *oneToOne[M, R]) LocalKey() Column {
	return oo.localKey
}

func (oo *oneToOne[M, R]) ForeignDef() Model {
	return oo.rel
}

func (oo *oneToOne[M, R]) Ref(m any) any {
	model, ok := m.(*M)
	if !ok {
		panic("invalid model type")
	}

	return oo.ref(model)
}

type oneToMany[M, R any] struct {
	self       TypedModel[M]
	rel        TypedModel[R]
	foreignKey Column
	localKey   Column
	ref        func(*M) *[]R
}

func (om *oneToMany[M, R]) Append(m any, r any) {
	model, ok := m.(*M)
	if !ok {
		panic("invalid model type")
	}

	rel, ok := r.(*R)
	if !ok {
		panic("invalid model type")
	}

	*om.ref(model) = append(*om.ref(model), *rel)
}

func (om *oneToMany[M, R]) ForeignKey() Column {
	return om.foreignKey
}

func (om *oneToMany[M, R]) LocalKey() Column {
	return om.localKey
}

func (om *oneToMany[M, R]) ForeignDef() Model {
	return om.rel
}

type OneToMany interface {
	Append(m any, r any)
	LocalKey() Column
	ForeignKey() Column
	ForeignDef() Model
}

type OneToOne interface {
	Set(any, any)
	LocalKey() Column
	ForeignKey() Column
	ForeignDef() Model
}

func BelongsTo[M, R any](localKey ModelCol[M], foreignKey ModelCol[R], ref func(m *M) *R) OneToOne {
	return &oneToOne[M, R]{
		self:       localKey.GetInnerModel(),
		rel:        foreignKey.GetInnerModel(),
		localKey:   localKey,
		foreignKey: foreignKey,
		ref:        ref,
	}
}

func HasOne[M, R any](localKey ModelCol[M], foreignKey ModelCol[R], ref func(m *M) *R) OneToOne {
	return &oneToOne[M, R]{
		self:       localKey.GetInnerModel(),
		rel:        foreignKey.GetInnerModel(),
		localKey:   localKey,
		foreignKey: foreignKey,
		ref:        ref,
	}
}

func HasMany[M, R any](localKey ModelCol[M], foreignKey ModelCol[R], ref func(m *M) *[]R) OneToMany {
	return &oneToMany[M, R]{
		self:       localKey.GetInnerModel(),
		rel:        foreignKey.GetInnerModel(),
		localKey:   localKey,
		foreignKey: foreignKey,
		ref:        func(m *M) *[]R { return ref(m) },
	}
}
