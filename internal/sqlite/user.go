package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal/apperror"
	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/mattn/go-sqlite3"
)

// UserStore handles user database operations.
type UserStore struct {
	db *DB
}

// NewUserStore returns a new UserStore.
func NewUserStore(db *DB) *UserStore {
	return &UserStore{db: db}
}

// GetUserByID returns a user from the database.
func (us *UserStore) GetUserByID(ctx context.Context, id string) (user.User, error) {
	sb := builder.NewSelectBuilder()
	sb.Select(
		"id",
		"name",
		"email",
		"created_at",
		"updated_at",
	)
	sb.From("user")
	sb.Where(sb.EQ("id", id))

	q, args := sb.Build()

	logger := log.FromContext(ctx)
	logger.Info("Get user by ID", "query", q, "args", args)

	var dest struct {
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		Email     string    `db:"email"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := us.db.Reader.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = apperror.New(apperror.ErrorCodeNotFound, "user not found")
		}
		return user.User{}, fmt.Errorf("sqlite: get user by id: %w", err)
	}
	return user.User(dest), nil
}

// CreateUser creates a new user in the database.
func (us *UserStore) CreateUser(ctx context.Context, input user.CreateUserInput) (user.User, error) {
	ib := builder.NewInsertBuilder()
	ib.InsertInto("user")
	ib.Cols(
		"id",
		"name",
		"email",
	)
	ib.Values(
		input.ID,
		input.Name,
		input.Email,
	)
	ib.Returning("created_at", "updated_at")

	q, args := ib.Build()

	logger := log.FromContext(ctx)
	logger.Info("Create user", "query", q, "args", args)

	var dest struct {
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := us.db.Writer.GetContext(ctx, &dest, q, args...); err != nil {
		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			err = apperror.New(apperror.ErrorCodeConflict, "email already taken")
		}
		return user.User{}, fmt.Errorf("sqlite: create user: %w", err)
	}

	return user.User{
		ID:        input.ID,
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: dest.CreatedAt,
		UpdatedAt: dest.UpdatedAt,
	}, nil
}
