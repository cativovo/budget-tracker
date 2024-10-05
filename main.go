package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pressly/goose/v3"
)

//go:embed internal/store/migrations/*.sql
var embedMigrations embed.FS

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DB, cfg.DB.SSL,
	)
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	goose.SetBaseFS(embedMigrations)
	db := stdlib.OpenDBFromPool(dbpool)
	if err := goose.Up(db, "internal/store/migrations"); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", hello(dbpool))

	log.Fatal(e.Start(":" + cfg.Port))
}

func hello(dbpool *pgxpool.Pool) func(echo.Context) error {
	return func(c echo.Context) error {
		rows, err := dbpool.Query(c.Request().Context(), "select * from account")
		if err != nil {
			return err
		}

		for rows.Next() {
			var id string
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				return err
			}

			fmt.Println(id, name)
		}

		return c.HTML(http.StatusOK, "hello world")
	}
}
