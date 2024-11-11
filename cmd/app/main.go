package main

import (
	"context"
	"log"

	budgettracker "github.com/cativovo/budget-tracker"
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

	v := vite.NewVite(vite.ViteConfig{
		IsDev:  cfg.Env == "development",
		DistFS: budgettracker.Dist,
	})

	server := server.NewServer(server.Resource{
		TransactionStore: queries,
		AssetsStore:      v,
	})

	log.Fatal(server.Start(":" + cfg.Port))
}
