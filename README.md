# MODB

MODB is an experimental ORM for GO which doesn't require code generation or reflection.

## Defining models

Models are defined declaratively:

```go 
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

var (
	UserModel = model.Define(model.ModelDefinition[User, userColumns, int]{
		Table: "users",
		Schema: userColumns{
			Id:   model.AutoIncrement("id", func(u *User) *int { return &u.Id }),
			Name: model.Col("name", func(u *User) *string { return &u.Name }),
			Age:  model.Col("age", func(u *User) *int { return &u.Age }),
		},
		PK: func(s userColumns) model.TypedCol[User, int] {
			return s.Id
		},
	})

	UserRelations = &userRelations{
		Posts: model.HasMany(UserModel.Cols().Id, PostModel.Cols().UserId, func(u *User) *[]Post { return &u.Posts }),
	}
)
```


## Querying
MODB use the repository pattern. From your model definition and your db (or a transaction) you can create a repository object. The repository is a generic object:

```go
    userRepo := repo.New(db, models.UserModel)
    postRepo := repo.New(db, models.PostModel)
```


Then you can call the repository methods. These methods are type-safe and inherit their types from the provided model:
```go
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

	for _, post := range user.Posts {
		println(post.Title)
	}
```

Output:
```
User:
Lucas
25
Posts:
Hello
Goodbye
```


> You can find a full example in the [example](example) folder.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details