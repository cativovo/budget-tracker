package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/log"
	mmiddleware "github.com/cativovo/budget-tracker/internal/server/middleware"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents an HTTP server.
type Server struct {
	router          *chi.Mux
	categoryService *category.Service
	userService     *user.Service
	addr            string
}

// New returns a new Server.
func New(
	us *user.Service,
	cs *category.Service,
) *Server {
	return &Server{
		router:          chi.NewRouter(),
		userService:     us,
		categoryService: cs,
	}
}

func (s *Server) mountRoutes() {
	// TODO: should we put middleware setup into a separate function?
	s.router.Use(middleware.RequestID)
	s.router.Use(mmiddleware.RequestLogger())
	s.router.Use(middleware.Heartbeat("/ping"))
	s.router.Use(middleware.Recoverer)

	mountDocsRoute(s.router)

	s.router.Group(func(r chi.Router) {
		api := huma.NewGroup(humachi.New(r, humaConfig), "/api")

		// TODO: update when doing the auth
		api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
			r, _ := humachi.Unwrap(ctx)
			u, err := s.userService.GetUserByID(r.Context(), "69")
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
				return
			}
			ctx = huma.WithContext(ctx, user.WithContext(r.Context(), u))
			next(ctx)
		})

		mountCategoryRoutes(api, s.categoryService)
	})
}

func (s *Server) setup() {
	s.mountRoutes()
}

// Serve starts the HTTP server and gracefully shuts it down on interrupt signals.
func (s *Server) Serve(ctx context.Context, addr string) error {
	s.setup()
	s.addr = addr

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	go func() {
		ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
		defer cancel()

		<-ctx.Done()

		slog.Info("Shutting down the server")
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown the server", log.ErrAttr(err))
		}
	}()

	slog.Info("Server listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server: listen and serve: %w", err)
	}
	return nil
}

func (s *Server) Addr() string {
	return s.addr
}
