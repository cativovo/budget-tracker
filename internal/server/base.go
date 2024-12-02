package server

import (
	"net/http"
	"time"

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

func getLogger(c echo.Context) *zap.SugaredLogger {
	return c.Get(ctxKeyLogger).(*zap.SugaredLogger)
}

func (s Server) Start(addr string) error {
	s.r.Logger.Infow("Starting server", "address", addr)
	return s.e.Start(addr)
}

func applyMiddlewares(e *echo.Echo, r Resource) {
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:      true,
		LogLatency:       true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		BeforeNextFunc: func(c echo.Context) {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = c.Response().Header().Get(echo.HeaderXRequestID)
			}

			logger := r.Logger.With(
				"request_id", requestID,
				"protocol", c.Request().Proto,
				"remote_ip", c.RealIP(),
				"host", c.Request().Host,
				"method", c.Request().Method,
				"uri", c.Request().RequestURI,
				"referer", c.Request().Referer(),
				"user_agent", c.Request().UserAgent(),
			)
			c.Set(ctxKeyLogger, logger)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger := getLogger(c)
			logger.Infow(
				"Request handled",
				"start_time", v.StartTime.Format(time.RFC3339),
				"latency_ms", v.Latency.Milliseconds(),
				"status", v.Status,
				"error", v.Error,
				"content_length", v.ContentLength,
				"response_size", v.ResponseSize,
			)
			return nil
		},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: http.FS(ui.DistDirFS),
	}))
}
