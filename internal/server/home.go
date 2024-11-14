package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cativovo/budget-tracker/internal/constants"
	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/cativovo/budget-tracker/internal/ui/pages"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type homeResource struct {
	transactionStore TransactionStore
	assetsStore      AssetsStore
}

func (hr homeResource) mountRoutes(e *echo.Echo) {
	e.GET("/", hr.homePage)
	e.GET("/test", hr.homePage)
}

func (hr homeResource) homePage(c echo.Context) error {
	type query struct {
		StartDate        string                       `query:"start_date"`
		EndDate          string                       `query:"end_date"`
		TransactionTypes []repository.TransactionType `query:"transaction_type"`
		Page             int                          `query:"page"`
	}

	var q query
	if err := c.Bind(&q); err != nil {
		// TODO: handle properly
		return err
	}

	// TODO: get from cookie
	accountID, err := uuid.Parse("c353143b-b608-4d63-9dc4-f7c434f2a3ff")
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	now := time.Now()
	year := now.Year()
	month := now.Month()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	if q.StartDate != "" && q.EndDate != "" {
		s, sErr := time.Parse(constants.DateFormat, q.StartDate)
		e, eErr := time.Parse(constants.DateFormat, q.EndDate)

		if sErr == nil && eErr == nil {
			startDate = s
			endDate = e
		}
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
		transactionTypes = []repository.TransactionType{
			repository.TransactionTypeIncome,
			repository.TransactionTypeExpense,
		}
	}

	queryParams := fmt.Sprintf(
		"page=%d&start_date=%s&end_date=%s",
		page+1,
		startDate.Format(constants.DateFormat),
		endDate.Format(constants.DateFormat),
	)
	for _, v := range transactionTypes {
		queryParams += fmt.Sprintf("&transaction_type=%s", v)
	}

	r, err := hr.transactionStore.ListTransactionsByDate(c.Request().Context(), repository.ListTransactionsByDateParams{
		TransactionTypes: transactionTypes,
		AccountID:        accountID,
		Limit:            limit,
		Offset:           offset,
		StartDate:        startDate,
		EndDate:          endDate,
	})
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return render(c, http.StatusOK, pages.Home(pages.HomeProps{
		Transactions: r.TransactionsByDate,
		QueryParams:  queryParams,
		HasNextPage:  false,
		AssetsStore:  hr.assetsStore,
	}))
}
