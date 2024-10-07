package server

import "github.com/labstack/echo/v4"

func isHxRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}
