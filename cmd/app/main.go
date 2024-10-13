package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/cativovo/budget-tracker/assets"
	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/server"
	"github.com/cativovo/budget-tracker/internal/store"
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
	assetsFs := getAssetFS(cfg.Env == "development")

	server := server.NewServer(server.Resource{
		TransactionStore: queries,
		AssetsFS:         assetsFs,
	})

	log.Fatal(server.Start(":" + cfg.Port))
}

func getAssetFS(useOS bool) http.FileSystem {
	if useOS {
		log.Println("using live mode")
		return http.FS(os.DirFS("assets"))
	}

	log.Println("using embed mode")
	return http.FS(assets.Assets)
}
