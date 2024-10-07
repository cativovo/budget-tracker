package server

import (
	"net/http"
	"time"

	"github.com/cativovo/budget-tracker/internal/store"
	"github.com/cativovo/budget-tracker/internal/ui/pages"
	"github.com/labstack/echo/v4"
)

type homeResource struct {
	expenseStore ExpenseStore
}

func (hr homeResource) mountRoutes(e *echo.Echo) {
	e.GET("/", hr.homePage)
}

func (hr homeResource) homePage(c echo.Context) error {
	// TODO: get from cookie
	accountID, err := store.NewUUID("52a4a56d-1ce0-4e77-92b9-e1051437ffee")
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	now := time.Now()
	year := now.Year()
	month := now.Month()
	monthStart := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	startDate, err := store.NewDate(monthStart)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	endDate, err := store.NewDate(monthEnd)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	expenses, err := hr.expenseStore.ListExpenses(c.Request().Context(), store.ListExpensesParams{
		AccountID: accountID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return err
	}

	return render(c, http.StatusOK, pages.Home(pages.HomeProps{
		Expenses: expenses,
	}))
}
