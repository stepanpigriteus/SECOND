package postgress

import (
	"context"
	"database/sql"
	"fmt"

	"1337b04rd/internal/domain/entity"
)

type PostgresSessionRepository struct {
	db *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) CreateSession(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	fmt.Printf(">>> Repository: Creating session with ID=%s, UserID=%s\n", session.ID, session.UserID)

	query := `INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, session.ID, session.UserID, session.CreatedAt, session.ExpiresAt).Scan(&session.ID)
	if err != nil {
		fmt.Printf(">>> Repository: Error creating session: %v\n", err)
		return nil, fmt.Errorf("не удалось создать сессию: %v", err)
	}

	fmt.Printf(">>> Repository: Session created successfully: ID=%s\n", session.ID)
	return session, nil
}

func (r *PostgresSessionRepository) GetSessionByID(ctx context.Context, id string) (*entity.Session, error) {
	fmt.Printf(">>> Repository: Getting session by ID=%s\n", id)

	query := `SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = $1`
	session := &entity.Session{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf(">>> Repository: Session not found: %s\n", id)
			return nil, fmt.Errorf("сессия с ID %s не найдена", id)
		}
		fmt.Printf(">>> Repository: Error querying session: %v\n", err)
		return nil, fmt.Errorf("не удалось получить сессию: %v", err)
	}

	fmt.Printf(">>> Repository: Session found: ID=%s, UserID=%s\n", session.ID, session.UserID)
	return session, nil
}

func (r *PostgresSessionRepository) DeleteSession(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}
