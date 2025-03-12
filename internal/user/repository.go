package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, u CreateUserReq) (User, error)
	DeleteUser(ctx context.Context, id string) error
}
