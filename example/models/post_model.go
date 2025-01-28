package models

import (
	"github.com/lucasjacques/modb/model"
)

type (
	Post struct {
		Id     int
		UserId int
		Title  string
		Body   string
		User   User
	}

	postSchema struct {
		Id     model.TypedCol[Post, int]
		UserId model.TypedCol[Post, int]
		Title  model.TypedCol[Post, string]
		Body   model.TypedCol[Post, string]
	}

	postRelations struct {
		User model.Relation
	}
)

var (
	PostModel = model.Define(model.ModelDefinition[Post, postSchema, int]{
		Table: "posts",
		Schema: postSchema{
			Id:     model.AutoIncrement("id", func(p *Post) *int { return &p.Id }),
			UserId: model.Col("user_id", func(p *Post) *int { return &p.UserId }),
			Title:  model.Col("title", func(p *Post) *string { return &p.Title }),
			Body:   model.Col("body", func(p *Post) *string { return &p.Body }),
		},
		PK: func(s postSchema) model.TypedCol[Post, int] {
			return s.Id
		},
	})

	PostRelations = &postRelations{
		User: model.BelongsTo(PostModel.Cols().UserId, UserModel.Cols().Id, func(p *Post) *User { return &p.User }),
	}
)
