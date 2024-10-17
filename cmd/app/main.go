package main

import (
	"context"
	"log"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/server"
	"github.com/cativovo/budget-tracker/internal/store"
	"github.com/cativovo/budget-tracker/internal/vite"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbpool, err := store.InitDB(context.Background(), cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	queries := store.New(dbpool)
	v := vite.NewVite(cfg.Env == "development")

	server := server.NewServer(server.Resource{
		TransactionStore: queries,
		AssetsStore:      v,
	})

	log.Fatal(server.Start(":" + cfg.Port))
}
