package server

import (
	"fmt"
	"math"
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
	accountID, err := store.NewUUID("c6d64bb9-2d0e-43c1-aa03-912912351f42")
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

	transactionTypes := q.TransactionTypes
	if len(transactionTypes) == 0 {
		transactionTypes = []int16{constants.TransactionTypeExpense, constants.TransactionTypeIncome}
	}

	transactionsWithCount, err := hr.transactionStore.ListTransactionsWithCount(c.Request().Context(), store.ListTransactionsParams{
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

	totalPages := int(math.Ceil(float64(transactionsWithCount.CountTotal) / float64(itemsPerPage)))
	hasNextPage := page < totalPages

	queryParams := fmt.Sprintf(
		"page=%d&start_date=%s&end_date=%s",
		page+1,
		start.Format(constants.DateFormat),
		end.Format(constants.DateFormat),
	)
	for _, v := range transactionTypes {
		queryParams += fmt.Sprintf("&transaction_type=%d", v)
	}

	return render(c, http.StatusOK, pages.Home(pages.HomeProps{
		TransactionsWithCount: transactionsWithCount,
		QueryParams:           queryParams,
		HasNextPage:           hasNextPage,
	}))
}
