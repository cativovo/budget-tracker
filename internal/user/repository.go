package user

import "context"

type Repository interface {
	UserByID(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, u CreateUserReq) (User, error)
	DeleteUser(ctx context.Context, id string) error
}
