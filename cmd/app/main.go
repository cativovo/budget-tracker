package main

import (
	"context"
	"log"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/cativovo/budget-tracker/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	r, err := repository.NewRepository(context.Background(), cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	server := server.NewServer(server.Resource{})

	log.Fatal(server.Start(":" + cfg.Port))
}
