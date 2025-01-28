package modb

import (
	"context"

	"github.com/lucasjacques/modb/model"
	"github.com/lucasjacques/modb/queries"
)

func (r *ModelRepository[M, PK, C]) getColumns(columns []model.Column) []string {
	cols := make([]string, 0, len(columns))
	for _, col := range columns {
		cols = append(cols, col.FQCN())
	}
	return cols
}

type queryOptions struct {
	limit int
	load  []model.Relation
	where queries.Expr
}

func (m *ModelRepository[M, PK, C]) query(ctx context.Context, opts *queryOptions) ([]M, error) {
	table, columns := m.model.GetTable(), m.getColumns(m.model.GetColumns())
	var joins []*queries.Join
	var eagerLoads []model.OneToOne
	var prefetches []model.OneToMany
	for _, rel := range opts.load {
		switch r := rel.(type) {
		case model.OneToOne:
			eagerLoads = append(eagerLoads, r)
		case model.OneToMany:
			prefetches = append(prefetches, r)
		}
	}

	for _, eagerLoad := range eagerLoads {
		eagerModel := eagerLoad.ForeignDef()
		eagerTable, eagerCols := eagerModel.GetTable(), m.getColumns(eagerModel.GetColumns())
		join := queries.InnerJoin(eagerTable).On(queries.EQ(eagerLoad.ForeignKey(), eagerLoad.LocalKey()))
		joins = append(joins, join)
		columns = append(columns, eagerCols...)
	}

	query := queries.NewQuery(m.db.NewParamsSet()).Select(table, columns...)
	for _, join := range joins {
		query = query.Join(join)
	}

	if opts.where != nil {
		query = query.Where(opts.where)
	}

	if opts.limit > 0 {
		query = query.Limit(opts.limit)
	}

	str, values := query.Build()

	rows, err := m.db.Query(ctx, str, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var models []M
	for rows.Next() {
		var dests []any
		modelDests := m.model.NewDests()
		dests = append(dests, modelDests...)

		var eagerDests [][]any
		for _, eagerLoad := range eagerLoads {
			dest := eagerLoad.ForeignDef().NewDests()
			dests = append(dests, dest...)
			eagerDests = append(eagerDests, dest)
		}

		err := rows.Scan(dests...)
		if err != nil {
			return nil, err
		}

		model, err := m.model.FromDestsTyped(modelDests)
		if err != nil {
			return nil, err
		}

		for i, eagerLoad := range eagerLoads {
			eagerModel, err := eagerLoad.ForeignDef().FromDests(eagerDests[i])
			if err != nil {
				return nil, err
			}

			eagerLoad.Set(model, eagerModel)
		}

		models = append(models, *model)
	}

	err = m.makePrefetches(ctx, models, prefetches)
	if err != nil {
		return nil, err
	}

	return models, nil
}

func multiValuesExpr(values []any) []queries.Expr {
	exprs := make([]queries.Expr, 0, len(values))
	for _, value := range values {
		exprs = append(exprs, queries.Value(value))
	}

	return exprs
}

func (m *ModelRepository[M, PK, S]) makePrefetches(ctx context.Context, results []M, prefetches []model.OneToMany) error {
	for _, prefetch := range prefetches {
		var invals []any
		for _, model := range results {
			value, err := prefetch.LocalKey().ValueFromModel(&model)
			if err != nil {
				return err
			}
			invals = append(invals, value)
		}

		prefetchQuery := queries.NewQuery(m.db.NewParamsSet()).
			Select(prefetch.ForeignDef().GetTable(), m.getColumns(prefetch.ForeignDef().GetColumns())...).
			Where(queries.IN(prefetch.ForeignKey(), multiValuesExpr(invals)))

		str, values := prefetchQuery.Build()

		rows, err := m.db.Query(ctx, str, values...)
		if err != nil {
			return err
		}

		defer rows.Close()

		var prefetchModels []any
		for rows.Next() {
			dests := prefetch.ForeignDef().NewDests()
			err := rows.Scan(dests...)
			if err != nil {
				return err
			}

			prefetchModel, err := prefetch.ForeignDef().FromDests(dests)
			if err != nil {
				return err
			}

			prefetchModels = append(prefetchModels, prefetchModel)
		}
		prefetchIndex := 0
		for i := range results {
			fk, err := prefetch.LocalKey().ValueFromModel(&results[i])
			if err != nil {
				return err
			}

			for prefetchIndex < len(prefetchModels) {
				prefetchModel := prefetchModels[prefetchIndex]
				prefetchFk, err := prefetch.ForeignKey().ValueFromModel(prefetchModel)
				if err != nil {
					return err
				}

				if fk == prefetchFk {
					prefetch.Append(&results[i], prefetchModel)
					prefetchIndex++
				} else {
					break
				}

			}

		}

	}

	return nil
}

func (m *ModelRepository[M, PK, S]) FindById(ctx context.Context, id PK, opts ...FindOpt) (*M, error) {
	findOps := buildFindOpts(opts)
	findOps.where = queries.EQ(m.model.PrimaryKey(), queries.Value(id))
	results, err := m.query(ctx, findOps)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	return &results[0], nil
}

func (m *ModelRepository[M, PK, S]) Find(ctx context.Context, opts ...FindOpt) ([]M, error) {
	findOps := buildFindOpts(opts)
	return m.query(ctx, findOps)
}

func (m *ModelRepository[M, PK, C]) FindOne(ctx context.Context, opts ...FindOpt) (*M, error) {
	findOps := buildFindOpts(opts)

	findOps.limit = 1

	results, err := m.query(ctx, findOps)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	return &results[0], nil
}
