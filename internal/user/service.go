package user

import "context"

// Service handles user operations.
type Service struct {
	userStore userStore
}

// NewService returns a new Service.
func NewService(us userStore) *Service {
	return &Service{userStore: us}
}

// GetUserByID returns a user by ID from the userStore.
func (s *Service) GetUserByID(ctx context.Context, id string) (User, error) {
	return s.userStore.GetUserByID(ctx, id)
}
