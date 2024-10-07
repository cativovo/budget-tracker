package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cativovo/budget-tracker/internal/constants"
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
	type query struct {
		StartDate string `query:"start_date"`
		EndDate   string `query:"end_date"`
		Page      int    `query:"page"`
	}

	var q query
	if err := c.Bind(&q); err != nil {
		// TODO: handle properly
		return err
	}

	// TODO: get from cookie
	accountID, err := store.NewUUID("52a4a56d-1ce0-4e77-92b9-e1051437ffee")
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	now := time.Now()
	year := now.Year()
	month := now.Month()
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	if q.StartDate != "" && q.EndDate != "" {
		s, sErr := time.Parse(constants.DateFormat, q.StartDate)
		e, eErr := time.Parse(constants.DateFormat, q.EndDate)

		if sErr == nil && eErr == nil {
			start = s
			end = e
		}
	}

	startDate, err := store.NewDate(start)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	endDate, err := store.NewDate(end)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	page := q.Page
	if page == 0 {
		page = 1
	}

	const itemsPerPage = 10
	offset := 0
	limit := page * itemsPerPage
	if isHxRequest(c) {
		offset = (page - 1) * itemsPerPage
		limit = itemsPerPage
	}

	expenses, err := hr.expenseStore.ListExpenses(c.Request().Context(), store.ListExpensesParams{
		AccountID: accountID,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return render(c, http.StatusOK, pages.Home(pages.HomeProps{
		Expenses: expenses,
		QueryParams: fmt.Sprintf(
			"page=%d&start_date=%s&end_date=%s",
			page+1,
			start.Format(constants.DateFormat),
			end.Format(constants.DateFormat),
		),
	}))
}
