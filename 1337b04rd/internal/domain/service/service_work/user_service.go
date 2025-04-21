package servicework

import (
	"context"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return nil
}
