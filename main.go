package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/store"
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

	queries := store.New(dbpool)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", hello1(queries))
	e.POST("/", hello2(queries))

	log.Fatal(e.Start(":" + cfg.Port))
}

type expenseStore interface {
	ListExpenses(ctx context.Context, arg store.ListExpensesParams) ([]store.ListExpensesRow, error)
	CreateExpense(ctx context.Context, arg store.CreateExpenseParams) (store.CreateExpenseRow, error)
}

func hello1(es expenseStore) func(echo.Context) error {
	return func(c echo.Context) error {
		accountID, err := store.NewUUID("d31d0b0a-f632-420d-a13e-ec093bef11b2")
		if err != nil {
			return err
		}

		startDate, err := store.NewDate("2024-09-01")
		if err != nil {
			return err
		}

		endDate, err := store.NewDate("2024-10-30")
		if err != nil {
			return err
		}

		expenses, err := es.ListExpenses(c.Request().Context(), store.ListExpensesParams{
			AccountID: accountID,
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, expenses)
	}
}

func hello2(es expenseStore) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		amount, err := store.NewNumeric(c.FormValue("amount"))
		if err != nil {
			return err
		}

		description, err := store.NewText(c.FormValue("description"))
		if err != nil {
			return err
		}

		date, err := store.NewDate(c.FormValue("date"))
		if err != nil {
			return err
		}

		categoryID, err := store.NewUUID(c.FormValue("category_id"))
		if err != nil {
			return err
		}

		accountID, err := store.NewUUID("d31d0b0a-f632-420d-a13e-ec093bef11b2")
		if err != nil {
			return err
		}

		expense, err := es.CreateExpense(c.Request().Context(), store.CreateExpenseParams{
			Name:        name,
			Amount:      amount,
			Description: description,
			Date:        date,
			CategoryID:  categoryID,
			AccountID:   accountID,
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, expense)
	}
}
