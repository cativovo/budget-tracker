package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func requestLogger(parentLogger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			rid := middleware.GetReqID(r.Context())

			startTime := time.Now()
			logger := parentLogger.With("request_id", rid)

			logger.Infow(
				"Processing Request",
				"protocol", r.Proto,
				"host", r.Host,
				"uri", r.RequestURI,
				"method", r.Method,
				"remote_ip", r.RemoteAddr,
				"referer", r.Referer(),
				"user_agent", r.UserAgent(),
				"start_time", startTime.Format(time.RFC3339),
			)

			defer func() {
				endTime := time.Now()
				latency := endTime.Sub(startTime)

				logger.Infow(
					"Request handled",
					"latency_ms", latency.Milliseconds(),
					"uri", r.RequestURI,
					"method", r.Method,
					"status", ww.Status(),
					"response_size", ww.BytesWritten(),
					"content_length", r.Header.Get(headerContentLength),
					"end_time", endTime.Format(time.RFC3339),
				)
			}()

			ctx := context.WithValue(r.Context(), ctxKeyLogger, logger)
			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}

// Adds X-Request-Id header to the response
func addRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := middleware.GetReqID(r.Context())
		w.Header().Set(middleware.RequestIDHeader, rid)
		next.ServeHTTP(w, r)
	})
}
