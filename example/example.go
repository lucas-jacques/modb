package main

import (
	"context"
	"database/sql"
	"fmt"
	"modbexample/models"

	"github.com/lucasjacques/modb/drivers/sqldriver"
	"github.com/lucasjacques/modb/repo"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE posts (id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT, body TEXT);")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	moDB := sqldriver.NewMODB(db, sqldriver.FQCNDoubleQuotes)

	userRepo := repo.New(moDB, models.UserModel)
	postRepo := repo.New(moDB, models.PostModel)

	err = userRepo.Insert(ctx, &models.User{
		Name: "Lucas",
		Age:  25,
	})
	if err != nil {
		panic(err)
	}

	err = postRepo.Insert(ctx, &models.Post{
		UserId: 1,
		Title:  "Hello",
		Body:   "World",
	})
	if err != nil {
		panic(err)
	}

	err = postRepo.Insert(ctx, &models.Post{
		UserId: 1,
		Title:  "Goodbye",
		Body:   "World",
	})
	if err != nil {
		panic(err)
	}

	user, err := userRepo.FindById(ctx, 1, repo.Preload(models.UserRelations.Posts))
	if err != nil {
		panic(err)
	}

	fmt.Println("User:")
	fmt.Println(user.Name)
	fmt.Println(user.Age)
	fmt.Println("Posts:")

	for _, post := range user.Posts {
		fmt.Println(post.Title)
	}

}
