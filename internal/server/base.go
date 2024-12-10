package server

import (
	"context"
	"net/http"

	"github.com/cativovo/budget-tracker/internal/repository"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Resource struct {
	Logger *zap.SugaredLogger
	//  make interface to make it easier to test
	Repository *repository.Repository
}

type Server struct {
	resource Resource
	router   *chi.Mux
}

type ctxKey string

const (
	ctxKeyLogger ctxKey = "logger"
)

func NewServer(r Resource) *Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(addRequestIDHeader)
	router.Use(middleware.RealIP)
	router.Use(requestLogger(r.Logger))
	router.Use(middleware.Compress(5, "text/html", "text/css", "text/javascript"))
	router.Use(middleware.Recoverer)

	router.Handle("/*", spaHandler())

	router.Route("/api", func(r chi.Router) {
		config := huma.DefaultConfig("My Api", "0.0.1")
		config.Servers = []*huma.Server{
			{URL: "/api"},
		}
		api := humachi.New(r, config)

		entryResource{}.mountRoutes(api)
	})

	return &Server{
		resource: r,
		router:   router,
	}
}

func (s Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func getLogger(ctx context.Context) *zap.SugaredLogger {
	return ctx.Value(ctxKeyLogger).(*zap.SugaredLogger)
}
