package repository

import "1337b04rd/internal/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int) error
}
