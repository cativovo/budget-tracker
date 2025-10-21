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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
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

	u, err := userService.GetUserByID(ctx, "69")
	if err != nil {
		logger.Error("Failed to get the user", slog.String("key", "69"), log.Error(err))
		return
	}

	logger.Info("User get by id", slog.Any("user", u))
}
