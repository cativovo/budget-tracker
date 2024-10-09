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
	transactionStore TransactionStore
}

func (hr homeResource) mountRoutes(e *echo.Echo) {
	e.GET("/", hr.homePage)
}

func (hr homeResource) homePage(c echo.Context) error {
	type query struct {
		StartDate        string  `query:"start_date"`
		EndDate          string  `query:"end_date"`
		TransactionTypes []int16 `query:"transaction_type"`
		Page             int     `query:"page"`
	}

	var q query
	if err := c.Bind(&q); err != nil {
		// TODO: handle properly
		return err
	}

	// TODO: get from cookie
	accountID, err := store.NewUUID("37bafcd9-8578-4d06-aaef-0bc3bd922d20")
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

	const itemsPerPage = 5
	offset := 0
	limit := page * itemsPerPage
	if isHxRequest(c) {
		offset = (page - 1) * itemsPerPage
		limit = itemsPerPage
	}

	transactionTypes := q.TransactionTypes
	if len(transactionTypes) == 0 {
		transactionTypes = []int16{constants.TransactionTypeExpense, constants.TransactionTypeIncome}
	}

	queryParams := fmt.Sprintf(
		"page=%d&start_date=%s&end_date=%s",
		page+1,
		start.Format(constants.DateFormat),
		end.Format(constants.DateFormat),
	)
	for _, v := range transactionTypes {
		queryParams += fmt.Sprintf("&transaction_type=%d", v)
	}

	r, err := hr.transactionStore.ListTransactionsByDate(c.Request().Context(), store.ListTransactionsByDateParams{
		TransactionTypes: transactionTypes,
		AccountID:        accountID,
		Limit:            int32(limit),
		Offset:           int32(offset),
		StartDate:        startDate,
		EndDate:          endDate,
	})
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	t, err := store.ParseListTransactionsByDateRows(r)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return render(c, http.StatusOK, pages.Home(pages.HomeProps{
		Transactions: t,
		QueryParams:  queryParams,
		HasNextPage:  false,
	}))
}
