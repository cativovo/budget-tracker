package server

import (
	"net/http"

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
	return render(c, http.StatusOK, pages.Home())
}
