package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/cativovo/budget-tracker/assets"
	"github.com/cativovo/budget-tracker/internal/ui/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Resource struct {
	ExpenseStore ExpenseStore
}

func NewServer(r Resource) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Gzip())

	homeResource{
		expenseStore: r.ExpenseStore,
	}.mountRoutes(e)

	assetHandler := http.FileServer(http.FS(assets.Assets))
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))

	e.HTTPErrorHandler = httpErrorHandler

	return e
}

func render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}

func httpErrorHandler(err error, c echo.Context) {
	statusCode := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		statusCode = he.Code
	}

	err = render(c, statusCode, pages.Error(pages.ErrorProps{
		StatusCode: statusCode,
	}))
	if err != nil {
		c.Logger().Error(err)
	}
}