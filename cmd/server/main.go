package main

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/cativovo/budget-tracker/internal/store"
	"github.com/cativovo/budget-tracker/internal/ui/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", index)
	e.GET("/hello1", hello1(queries))
	e.POST("/hello2", hello2(queries))

	log.Fatal(e.Start(":" + cfg.Port))
}

type expenseStore interface {
	ListExpenses(ctx context.Context, arg store.ListExpensesParams) ([]store.ListExpensesRow, error)
	CreateExpense(ctx context.Context, arg store.CreateExpenseParams) (store.CreateExpenseRow, error)
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func index(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Index())
}

func hello1(es expenseStore) func(echo.Context) error {
	return func(c echo.Context) error {
		accountID, err := store.NewUUID("52a4a56d-1ce0-4e77-92b9-e1051437ffee")
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
