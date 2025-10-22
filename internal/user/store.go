package user

import "context"

type userStore interface {
	GetUserByID(ctx context.Context, id string) (User, error)
	CreateUser(ctx context.Context, input CreateUserInput) (User, error)
	// TODO: add Update, Delete
	// Add test data for now, we'll need this later for authentication
}
