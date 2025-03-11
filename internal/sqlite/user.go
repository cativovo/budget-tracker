package sqlite

import (
	"context"

	"github.com/cativovo/budget-tracker/internal"
)

type UserRepository struct {
	db *DB
}

var _ internal.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, u internal.User) (internal.User, error) {
	panic("not yet implemented")
}
