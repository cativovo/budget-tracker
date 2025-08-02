package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/apperror"
	"github.com/cativovo/budget-tracker/internal/log"
	"github.com/cativovo/budget-tracker/internal/user"
	"github.com/huandu/go-sqlbuilder"
	"github.com/mattn/go-sqlite3"
)

// UserStore is a user store that uses a SQLite database.
type UserStore struct {
	db *DB
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// GetUser gets a user by id.
func (us *UserStore) GetUser(ctx context.Context, id string) (user.User, error) {
	logger := log.FromContext(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "email", "created_at", "updated_at")
	sb.From("user")
	sb.Where(sb.EQ("id", id))

	q, args := sb.Build()
	logger.Infow("Get user", "query", q, "args", args)

	var dest userDest
	if err := us.db.Reader.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = apperror.New(apperror.ErrorCodeNotFound, "user not found")
		}
		return user.User{}, fmt.Errorf("sqlite: get user: %w", err)
	}

	return user.User(dest), nil
}

// CreateUser creates a new user.
func (us *UserStore) CreateUser(ctx context.Context, input user.UserCreate) (user.User, error) {
	logger := log.FromContext(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("user")
	ib.Cols("id", "name", "email")
	ib.Values(input.ID, input.Name, input.Email)
	ib.Returning("*")

	q, args := ib.Build()
	logger.Infow("Insert user", "query", q, "args", args)

	var dest userDest
	if err := us.db.Writer.GetContext(ctx, &dest, q, args...); err != nil {
		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			err = apperror.New(apperror.ErrorCodeConflict, "email already taken")
		}
		return user.User{}, fmt.Errorf("sqlite: insert user: %w", err)
	}

	return user.User(dest), nil
}

// UpdateUser updates a user.
func (us *UserStore) UpdateUser(ctx context.Context, input user.UserUpdate) (user.User, error) {
	panic("err")
}

type userDest struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
