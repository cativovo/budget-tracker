package user

import "context"

type userStore interface {
	GetUser(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, user UserCreate) (User, error)
	UpdateUser(ctx context.Context, user UserUpdate) (User, error)
}
