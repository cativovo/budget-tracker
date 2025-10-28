package server_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cativovo/budget-tracker/internal/category"
	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/cativovo/budget-tracker/internal/server"
	"github.com/cativovo/budget-tracker/internal/sqlite"
	"github.com/cativovo/budget-tracker/internal/testutil"
	"github.com/cativovo/budget-tracker/internal/user"
)

var testServer *server.Server

func TestMain(m *testing.M) {
	db, c := newTestDB()
	defer c()

	userStore := sqlite.NewUserStore(db)
	userService := user.NewService(userStore)

	// TODO: remove this once auth is implemented
	_, err := userService.CreateUser(context.Background(), user.CreateUserInput{
		ID:    "69",
		Name:  "Juan Usa",
		Email: "juanusa@email.com",
	})
	if err != nil {
		slog.Error("Failed to create a test user", log.ErrAttr(err))
		os.Exit(1)
	}

	categoryStore := sqlite.NewCategoryStore(db)
	categoryService := category.NewService(categoryStore)

	testServer = server.New(userService, categoryService)
	go func() {
		// TODO: make the addr configurable
		if err := testServer.Serve(context.Background(), ":6969"); err != nil {
			slog.Error("Failed to start the server", log.ErrAttr(err))
			os.Exit(1)
		}
	}()

	waitForServer(testServer)

	slog.Info("Disabling the logs")
	testutil.DisableLogs()

	code := m.Run()
	os.Exit(code)
}

func waitForServer(s *server.Server) {
	slog.Info("Waiting for the server to be ready")

	interval := time.Millisecond * 100
	for {
		<-time.Tick(interval)

		resp, err := http.DefaultClient.Get("http://127.0.0.1" + s.Addr() + "/ping")
		if err != nil {
			slog.Error("Failed to make a ping request", log.ErrAttr(err))
			continue
		}
		if resp.StatusCode == http.StatusOK {
			slog.Info("Server is ready")
			return
		}

		slog.Info("Server is not yet ready", "status_code", resp.StatusCode)
	}
}

func newTestDB() (*sqlite.DB, func()) {
	tmp, err := os.CreateTemp("", "")
	if err != nil {
		slog.Error("Failed to create a temp file", log.ErrAttr(err))
		os.Exit(1)
	}

	defer func() {
		err := tmp.Close()
		if err != nil {
			slog.Error("Failed to close the temp file", log.ErrAttr(err))
			os.Exit(1)
		}
	}()

	db, err := sqlite.NewDB(context.Background(), tmp.Name())
	if err != nil {
		slog.Error("Failed to create new test db", log.ErrAttr(err))
		os.Exit(1)
	}

	err = db.Migrate(context.Background())
	if err != nil {
		slog.Error("Database migration failed", log.ErrAttr(err))
		os.Exit(1)
	}

	c := func() {
		n := tmp.Name()

		err := os.Remove(n)
		if err != nil {
			slog.Error("Failed to remove tmp file", log.ErrAttr(err))
		}

		err = os.Remove(fmt.Sprintf("%s-shm", n))
		if err != nil {
			slog.Error("Failed to remove shm tmp file", log.ErrAttr(err))
		}

		err = os.Remove(fmt.Sprintf("%s-wal", n))
		if err != nil {
			slog.Error("Failed to remove wal tmp file", log.ErrAttr(err))
		}
	}

	return db, c
}
