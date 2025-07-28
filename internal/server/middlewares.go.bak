package server

import (
	"context"
	"encoding/json"
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

// Sets X-Request-Id header to the response
func setResponseRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := middleware.GetReqID(r.Context())
		w.Header().Set(middleware.RequestIDHeader, rid)
		next.ServeHTTP(w, r)
	})
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				logger := getLogger(r.Context())
				logger.Errorw("Recovered from panic", "runtime_error", rvr)

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
					res := map[string]string{
						"message": "Internal server error",
					}
					if err := json.NewEncoder(w).Encode(res); err != nil {
						panic(err)
					}
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
