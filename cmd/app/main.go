package main

import (
	"context"
	"io/fs"
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
	viteManifest, err := budgettracker.Dist.Open("dist/.vite/manifest.json")
	if err != nil {
		panic(err)
	}
	defer viteManifest.Close()
	assets, err := fs.Sub(budgettracker.Dist, "dist/assets")
	if err != nil {
		panic(err)
	}
	v := vite.NewVite(cfg.Env == "development", viteManifest, assets)

	server := server.NewServer(server.Resource{
		TransactionStore: queries,
		AssetsStore:      v,
	})

	log.Fatal(server.Start(":" + cfg.Port))
}
