package server

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Resource struct {
	ExpenseStore ExpenseStore
}

func NewServer(r Resource) *echo.Echo {
	e := echo.New()

	homeResource{
		expenseStore: r.ExpenseStore,
	}.mountRoutes(e)

	return e
}

func render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
