package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/cativovo/budget-tracker/internal/server"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/user"
)

func main() {
	ho := &slog.HandlerOptions{
		AddSource: true,
	}
	// TODO: make configurable
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, ho))
	logger := slog.New(slog.NewTextHandler(os.Stdout, ho))
	slog.SetDefault(logger)

	ctx := context.Background()
	db, err := sqlite.NewDB(ctx, "budget_store.db")
	if err != nil {
		logger.Error("Failed to initialize the db", log.ErrAttr(err))
		return
	}

	err = db.Migrate(ctx)
	if err != nil {
		logger.Error("Failed to migrate db", log.ErrAttr(err))
		return
	}

	userStore := sqlite.NewUserStore(db)
	userService := user.NewService(userStore)

	// TODO: remove me once the auth is implemented
	_, err = userService.GetUserByID(ctx, "69")
	if err != nil {
		logger.Error("User doesn't exists, creating the user")

		_, err := userService.CreateUser(ctx, user.CreateUserInput{
			ID:    "69",
			Name:  "Juan Usa",
			Email: "juanusa@email.com",
		})
		if err != nil {
			logger.Error("Failed to create a user", log.ErrAttr(err))
			return
		}
	}

	categoryStore := sqlite.NewCategoryStore(db)
	categoryService := category.NewService(categoryStore)

	srv := server.New(
		userService,
		categoryService,
	)
	// TODO: make addr configurable
	if err := srv.Serve(ctx, ":4000"); err != nil {
		slog.Error("Failed to start the server", log.ErrAttr(err))
	}
}
