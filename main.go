package main

import (
	"log"
	"net/http"

	"github.com/cativovo/budget-tracker/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", hello)

	log.Fatal(e.Start(":" + cfg.Port))
}

func hello(c echo.Context) error {
	c.Logger().Warn("test")
	c.Logger().Error("test")
	return c.HTML(http.StatusOK, "hello world")
}
