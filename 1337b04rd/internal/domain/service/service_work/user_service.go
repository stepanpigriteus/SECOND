package servicework

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
	"context"
	"fmt"
	"time"
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
	if user.CharacterName == "" {
		return entity.User{}, fmt.Errorf("character name is required")
	}
	createdAt := time.Now()
	user.CreatedAt = createdAt

	createdUser, err := s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return entity.User{}, err
	}
	return *createdUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	fmt.Println(">>> GetUserById service called", id)

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return *user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	fmt.Println(">>> UpdateUser service called")
	updatedUser, err := s.userRepo.UpdateUser(ctx, &user)
	if err != nil {
		return entity.User{}, err
	}
	return *updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	fmt.Println(">>> DeleteUser service called")
	return s.userRepo.DeleteUser(ctx, id)
}
