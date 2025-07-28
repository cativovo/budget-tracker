package main

import (
	"context"

	"github.com/cativovo/budget-tracker/internal/service"
	"github.com/cativovo/budget-tracker/internal/sqlite"
)

func main() {
	db, err := sqlite.NewDB(context.TODO(), "/tmp/db.foo")
	if err != nil {
		panic(err)
	}
	userStore := sqlite.NewUserStore(db)
	userService := service.NewUserService(userStore)
	_ = userService
}
