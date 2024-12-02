package server

import (
	"net/http"

	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/cativovo/budget-tracker/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Resource struct {
	Logger *zap.SugaredLogger
	//  make interface to make it easier to test
	Repository *repository.Repository
}

type Server struct {
	e *echo.Echo
	r Resource
}

const (
	ctxKeyLogger = "logger"
)

func NewServer(r Resource) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	applyMiddlewares(e, r)

	api := e.Group("/api")
	api.GET("/foo", func(c echo.Context) error {
		t := []string{"uno", "dos", "tres"}
		return c.JSON(http.StatusOK, t)
	})

	type payload struct {
		Name string `json:"name"`
	}
	api.POST("/foo", func(c echo.Context) error {
		logger := getLogger(c)

		var p payload
		if err := c.Bind(&p); err != nil {
			logger.Error(err)
			return err
		}

		logger.Infow("Processing payload", "payload", p)

		acc, err := r.Repository.CreateAccount(c.Request().Context(), logger, repository.CreateAccountParams{
			Name: p.Name,
		})
		if err != nil {
			logger.Error(err)
			return err
		}

		logger.Infow("Created account successfully", "account", acc)

		return c.JSON(http.StatusOK, p)
	})

	return &Server{
		e: e,
		r: r,
	}
}

func (s Server) Start(addr string) error {
	s.r.Logger.Infow("Starting server", "address", addr)
	return s.e.Start(addr)
}

func getLogger(c echo.Context) *zap.SugaredLogger {
	return c.Get(ctxKeyLogger).(*zap.SugaredLogger)
}

func applyMiddlewares(e *echo.Echo, r Resource) {
	e.Use(middleware.RequestID())
	e.Use(RequestLogger(r.Logger))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: http.FS(ui.DistDirFS),
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger := getLogger(c)
			logger.Errorw("Runtime error", "error", err, "stack", string(stack))
			return err
		},
	}))
	e.Use(middleware.Gzip())
}
