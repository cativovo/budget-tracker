package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cativovo/budget-tracker/internal"
	"github.com/cativovo/budget-tracker/internal/ctxvalue"
	"github.com/cativovo/budget-tracker/internal/model"
	"github.com/huandu/go-sqlbuilder"
	"github.com/mattn/go-sqlite3"
)

type UserStore struct {
	db *DB
}

func NewUserStore(db *DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetUser(ctx context.Context, id string) (model.User, error) {
	logger := ctxvalue.LoggerFromContext(ctx)

	sb := sqlbuilder.SQLite.NewSelectBuilder()
	sb.Select("id", "name", "email", "created_at", "updated_at")
	sb.From("user")
	sb.Where(sb.EQ("id", id))

	q, args := sb.Build()
	logger.Infow("Get user", "query", q, "args", args)

	var dest struct {
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		Email     string    `db:"email"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := us.db.Reader.GetContext(ctx, &dest, q, args...); err != nil {
		if err == sql.ErrNoRows {
			err = internal.NewError(internal.ErrorCodeNotFound, "user not found")
		}
		return model.User{}, fmt.Errorf("sqlite: get user: %w", err)
	}

	return model.User(dest), nil
}

func (us *UserStore) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	logger := ctxvalue.LoggerFromContext(ctx)

	ib := sqlbuilder.SQLite.NewInsertBuilder()
	ib.InsertInto("user")
	ib.Cols("id", "name", "email")
	ib.Values(user.ID, user.Name, user.Email)
	ib.Returning("created_at", "updated_at")

	q, args := ib.Build()
	logger.Infow("Insert user", "query", q, "args", args)

	var dest struct {
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	if err := us.db.Writer.GetContext(ctx, &dest, q, args...); err != nil {
		var sqErr sqlite3.Error
		if errors.As(err, &sqErr) && sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			err = internal.NewError(internal.ErrorCodeConflict, "email already taken")
		}
		return model.User{}, fmt.Errorf("sqlite: insert user: %w", err)
	}

	user.CreatedAt = dest.CreatedAt
	user.UpdatedAt = dest.UpdatedAt

	return user, nil
}

func (us *UserStore) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	panic("err")
}
