package repository

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(user *entity.User) (*entity.User, error) {
	return user, nil
}

func (r *PostgresRepository) GetUserByID(id int) (*entity.User, error) {
	user := entity.User{}
	return &user, nil
}

func (r *PostgresRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	return user, nil
}

func (r *PostgresRepository) DeleteUser(id int) error {
	return nil
}
