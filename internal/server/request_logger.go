package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLogger(parentLogger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = res.Header().Get(echo.HeaderXRequestID)
			}

			startTime := time.Now()
			logger := parentLogger.With("request_id", requestID)
			c.Set(ctxKeyLogger, logger)

			logger.Infow(
				"Processing Request",
				"protocol", c.Request().Proto,
				"host", c.Request().Host,
				"uri", c.Request().RequestURI,
				"method", c.Request().Method,
				"remote_ip", c.RealIP(),
				"referer", c.Request().Referer(),
				"user_agent", c.Request().UserAgent(),
				"start_time", startTime.Format(time.RFC3339),
			)

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}

			endTime := time.Now()
			latency := endTime.Sub(startTime)

			logger.Infow(
				"Request handled",
				"latency_ms", latency.Milliseconds(),
				"uri", c.Request().RequestURI,
				"method", c.Request().Method,
				"status", res.Status,
				"response_size", res.Size,
				"content_length", req.Header.Get(echo.HeaderContentLength),
				"error", err,
				"end_time", endTime.Format(time.RFC3339),
			)

			return err
		}
	}
}
