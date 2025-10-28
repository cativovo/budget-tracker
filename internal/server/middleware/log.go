package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/go-chi/chi/v5/middleware"
)

// RequestLogger attaches a logger with context information to the request context.
func RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// NOTE: if we want to generate our own request ID, we must use middleware.RequestID from chi
			rid := middleware.GetReqID(r.Context())

			startTime := time.Now()
			logger := slog.With("request_id", rid)

			logger.Info(
				"Processing Request",
				"protocol", r.Proto,
				"host", r.Host,
				"uri", r.RequestURI,
				"method", r.Method,
				"remote_ip", r.RemoteAddr,
				"referer", r.Referer(),
				"user_agent", r.UserAgent(),
				"start_time", startTime,
			)

			defer func() {
				endTime := time.Now()
				latency := endTime.Sub(startTime)

				logger.Info(
					"Request handled",
					"latency_ms", latency.Milliseconds(),
					"uri", r.RequestURI,
					"method", r.Method,
					"status", ww.Status(),
					"response_size", ww.BytesWritten(),
					// "content_length", r.Header.Get(headerContentLength),
					"end_time", endTime,
				)
			}()

			ctx := log.WithContext(r.Context(), logger)
			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}
