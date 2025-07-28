package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLogger(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core).Sugar()
	router.Use(middleware.RequestID)
	router.Use(requestLogger(logger))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, getLogger(r.Context()))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "hi")
	})

	startTime := time.Now()
	router.ServeHTTP(w, r)
	endTime := time.Now()

	assert.Equal(t, 2, logs.Len())

	reqLogs := logs.All()[0]
	assert.Equal(t, "Processing Request", reqLogs.Message)
	reqContextMap := reqLogs.ContextMap()
	reqID := reqContextMap["request_id"].(string)
	assert.NotEmpty(t, reqID)
	assert.Equal(t, "HTTP/1.1", reqContextMap["protocol"])
	assert.Equal(t, "example.com", reqContextMap["host"])
	assert.Equal(t, "/", reqContextMap["uri"])
	assert.Equal(t, "GET", reqContextMap["method"])
	assert.Contains(t, reqContextMap, "referer")
	assert.Contains(t, reqContextMap, "user_agent")
	assert.NotEmpty(t, reqContextMap["remote_ip"])
	reqStartTime, err := time.Parse(time.RFC3339, reqContextMap["start_time"].(string))
	assert.Nil(t, err)
	assert.WithinDuration(t, startTime, reqStartTime, time.Second)

	resLogs := logs.All()[1]
	assert.Equal(t, "Request handled", resLogs.Message)
	resContextMap := resLogs.ContextMap()
	resEndTime, err := time.Parse(time.RFC3339, resContextMap["end_time"].(string))
	assert.Nil(t, err)
	assert.WithinDuration(t, endTime, resEndTime, time.Second)
	assert.Equal(t, resEndTime.Sub(reqStartTime).Milliseconds(), resContextMap["latency_ms"])
	assert.Equal(t, reqID, resContextMap["request_id"])
	assert.Equal(t, "/", resContextMap["uri"])
	assert.Equal(t, "GET", resContextMap["method"])
	assert.Equal(t, http.StatusOK, int(resContextMap["status"].(int64)))
	assert.Equal(t, 2, int(resContextMap["response_size"].(int64)))
	assert.Contains(t, resContextMap, "content_length")
}

func TestSetResponseRequestID(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(setResponseRequestID)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, middleware.GetReqID(r.Context()))
	})

	router.ServeHTTP(w, r)

	assert.NotEmpty(t, w.Header().Get(headerXRequestID))
	assert.Equal(t, w.Body.String(), w.Header().Get(headerXRequestID))
}

func TestRecoverer(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router := chi.NewRouter()
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core).Sugar()
	router.Use(requestLogger(logger))
	router.Use(recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		panic("don't panic, it's organic")
	})

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Equal(t, `{"message":"Internal server error"}`+"\n", w.Body.String())
	assert.Equal(t, "Recovered from panic", logs.All()[1].Message)
	assert.Contains(t, logs.All()[1].ContextMap()["runtime_error"], "don't panic, it's organic")
}
