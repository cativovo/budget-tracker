package server

import (
	"net/http"

	"github.com/cativovo/budget-tracker/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Resource struct{}

func NewServer(r Resource) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: http.FS(ui.DistDirFS),
	}))

	api := e.Group("/api")
	api.GET("/foo", func(c echo.Context) error {
		t := []string{"uno", "dos", "tres"}
		return c.JSON(http.StatusOK, t)
	})

	return e
}
