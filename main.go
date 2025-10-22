package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cativovo/budget-tracker/internal/log"
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
		logger.Error("Failed to initialize the db", log.Error(err))
		return
	}

	err = db.Migrate(ctx)
	if err != nil {
		logger.Error("Failed to migrate db", log.Error(err))
		return
	}

	userStore := sqlite.NewUserStore(db)
	userService := user.NewService(userStore)

	gotUser, err := userService.GetUserByID(ctx, "69")
	if err != nil {
		logger.Error("User doesn't exists, creating the user")

		gotUser, err = userService.CreateUser(ctx, user.CreateUserInput{
			ID:    "69",
			Name:  "Juan Usa",
			Email: "juanusa@email.com",
		})
		if err != nil {
			logger.Error("Failed to create a user", log.Error(err))
			return
		}
	}

	logger.Info("User get by id", slog.Any("user", gotUser))
}
