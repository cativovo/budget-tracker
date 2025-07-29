package main

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/user"
)

func main() {
	db, err := sqlite.NewDB(context.TODO(), "/tmp/db.foo")
	if err != nil {
		panic(err)
	}
	userStore := sqlite.NewUserStore(db)
	userService := user.NewService(userStore)
	_ = userService
}
