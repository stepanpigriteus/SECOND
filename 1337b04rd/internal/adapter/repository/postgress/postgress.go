package postgress

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
	"database/sql"
	"log"
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

func ConnectDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Connected to Database successfully!")

	return db, nil
}
