package server

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
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

			// Request
			reqBody := []byte{}
			if req.Body != nil {
				reqBody, _ = io.ReadAll(c.Request().Body)
			}
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			logger.Infow(
				"Processing Request",
				"protocol", c.Request().Proto,
				"host", c.Request().Host,
				"uri", c.Request().RequestURI,
				"method", c.Request().Method,
				"remote_ip", c.RealIP(),
				"referer", c.Request().Referer(),
				"user_agent", c.Request().UserAgent(),
				"request_body", string(reqBody),
				"start_time", startTime.Format(time.RFC3339),
			)

			// Response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

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
				"respose_body", resBody.String(),
				"content_length", req.Header.Get(echo.HeaderContentLength),
				"response_size", res.Size,
				"error", err,
				"end_time", endTime.Format(time.RFC3339),
			)

			return err
		}
	}
}

// https://github.com/labstack/echo/blob/fe2627778114fc774a1b10920e1cd55fdd97cf00/middleware/body_dump.go#L30
type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	err := http.NewResponseController(w.ResponseWriter).Flush()
	if err != nil && errors.Is(err, http.ErrNotSupported) {
		panic(errors.New("response writer flushing is not supported"))
	}
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

func (w *bodyDumpResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
