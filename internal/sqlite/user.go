package sqlite

import (
	"context"
	"time"

	"github.com/cativovo/budget-tracker/internal/user"
)

type UserStore struct {
	db *DB
}

func NewUserStore(db *DB) *UserStore {
	return &UserStore{db: db}
}

// TODO: get from db
func (us *UserStore) GetUserByID(ctx context.Context, id string) (user.User, error) {
	d := time.Date(2025, time.October, 21, 0, 0, 0, 0, time.UTC)
	return user.User{
		ID:        "1",
		Name:      "Juan Usa",
		Email:     "juanusa@email.com",
		CreatedAt: d,
		UpdatedAt: d,
	}, nil
}
