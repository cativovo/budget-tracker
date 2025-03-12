package user

import "context"

type UserRepository interface {
	FindUserByID(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, u CreateUserReq) (User, error)
	DeleteUser(ctx context.Context, id string) error
}
