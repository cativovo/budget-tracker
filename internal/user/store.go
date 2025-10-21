package user

import "context"

type userStore interface {
	GetUserByID(ctx context.Context, id string) (User, error)
	// TODO: add Create, Update, Delete
	// Add test data for now, we'll need this later for authentication
}
