package internal

import "context"

type User struct {
	ID    string
	Name  string
	Email string
}

type UserRepository interface {
	CreateUser(ctx context.Context, u User) (User, error)
}
