package main

import (
	"context"
	"log"

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

	server := server.NewServer(server.Resource{
		ExpenseStore: queries,
	})

	log.Fatal(server.Start(":" + cfg.Port))
}
