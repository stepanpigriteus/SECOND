package postgress

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// Выполняем запрос на вставку в таблицу users
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (character_name, avatar_url) VALUES ($1, $2) RETURNING id`,
		user.CharacterName, user.AvatarURL,
	).Scan(&user.ID)
	// Если ошибка, возвращаем её
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Возвращаем пользователя с ID
	return user, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User

	query := `SELECT id, avatar_url, character_name, custom_name, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&user.ID, &user.AvatarURL, &user.CharacterName, &user.CustomName, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return user, nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, id int) error {
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
