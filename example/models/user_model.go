package models

import (
	"github.com/lucasjacques/modb/model"
)

type (
	User struct {
		Id    int
		Name  string
		Age   int
		Posts []Post
	}

	userColumns struct {
		Id   model.TypedCol[User, int]
		Name model.TypedCol[User, string]
		Age  model.TypedCol[User, int]
	}

	userRelations struct {
		Posts model.Relation
	}
)

func (u *userColumns) PrimaryKey() model.TypedCol[User, int] {
	return u.Id
}

func (s *userColumns) Cols() []model.ModelCol[User] {
	return []model.ModelCol[User]{
		s.Id,
		s.Name,
		s.Age,
	}
}

var (
	UserModel = model.New(
		"users",
		&userColumns{
			Id:   model.AutoIncrement("id", func(u *User) *int { return &u.Id }),
			Name: model.Col("name", func(u *User) *string { return &u.Name }),
			Age:  model.Col("age", func(u *User) *int { return &u.Age }),
		},
	)

	UserRelations = &userRelations{
		Posts: model.HasMany(UserModel.Cols().Id, PostModel.Cols().UserId, func(u *User) *[]Post { return &u.Posts }),
	}
)
