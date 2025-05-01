package service

import (
	"context"

	"a1337b04rd/internal/domain/entity"
)

type UserService interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}
