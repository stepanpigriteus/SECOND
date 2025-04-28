package postgress

import (
	"context"
	"database/sql"

	"1337b04rd/internal/domain/entity"
)

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	return comment, nil
}

func (r *PostgresCommentRepository) GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, user_id, content, created_at
          FROM comments
         WHERE post_id = $1
         ORDER BY created_at ASC
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*entity.Comment
	for rows.Next() {
		var c entity.Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (r *PostgresCommentRepository) DeleteComment(ctx context.Context, id int) error {
	return nil
}
