package main

import (
	"fmt"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/cativovo/budget-tracker/internal/server"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		panic(err)
	}

	logger.Infow("Config details", "config", cfg)

	r, err := repository.NewRepository(cfg.DBPath)
	if err != nil {
		logger.Fatal(err)
	}
	defer r.Close()

	if err := r.Migrate(logger); err != nil {
		logger.Fatal(err)
	}

	s := server.NewServer(server.Resource{
		Logger:     logger,
		Repository: r,
	})

	logger.Fatal(s.Start(fmt.Sprintf(":%s", cfg.Port)))
}
