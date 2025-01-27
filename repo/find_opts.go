package repo

import (
	"github.com/lucasjacques/modb/model"
	"github.com/lucasjacques/modb/queries"
)

func buildFindOpts(opts []FindOpt) *queryOptions {
	var findOpts queryOptions
	for _, opt := range opts {
		opt(&findOpts)
	}
	return &findOpts
}

type FindOpt func(*queryOptions)

func Preload(r model.Relation) FindOpt {
	return func(opts *queryOptions) {
		opts.load = append(opts.load, r)
	}
}

func Limit(limit int) FindOpt {
	return func(opts *queryOptions) {
		opts.limit = limit
	}
}

func Where(expr queries.Expr) FindOpt {
	return func(opts *queryOptions) {
		opts.where = expr
	}
}
