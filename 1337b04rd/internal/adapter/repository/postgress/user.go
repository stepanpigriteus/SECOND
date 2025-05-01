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
    query := `SELECT id, character_name, avatar_url, custom_name FROM users WHERE id = $1`
    user := &entity.User{}
    
    // Используем указатели на sql.NullString для полей, которые могут быть NULL
    var customName sql.NullString
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.CharacterName,
        &user.AvatarURL,
        &customName, // Используем sql.NullString вместо &user.CustomName напрямую
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("пользователь с ID %d не найден", id)
        }
        return nil, fmt.Errorf("ошибка при получении пользователя: %v", err)
    }
    
    // Присваиваем значение только если оно не NULL
    if customName.Valid {
        user.CustomName = customName.String
    } else {
        user.CustomName = "" // или любое другое значение по умолчанию
    }
    
    return user, nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users
		SET custom_name = $1
		WHERE id = $2
		RETURNING id,character_name, custom_name
	`
	err := r.db.QueryRowContext(ctx, query, user.CustomName, user.ID).
		Scan(&user.ID, &user.CharacterName, &user.CustomName)
	if err != nil {
		return nil, err
	}

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
