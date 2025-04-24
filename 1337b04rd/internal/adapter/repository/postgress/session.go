package postgress

import (
	"1337b04rd/internal/domain/entity"
	"context"
	"database/sql"
	"fmt"
)

type PostgresSessionRepository struct {
	db *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (s *PostgresSessionRepository) CreateSession(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	return session, nil
}

func (r *PostgresSessionRepository) GetSessionByID(ctx context.Context, id string) (*entity.Session, error) {
	query := `SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = $1`
	session := &entity.Session{}
	err := r.db.QueryRow(query, id).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("сессия с ID %s не найдена", id)
		}
		return nil, fmt.Errorf("не удалось получить сессию: %v", err)
	}
	return session, nil
}

func (s *PostgresSessionRepository) DeleteSession(ctx context.Context, id string) error {
	return nil
}
