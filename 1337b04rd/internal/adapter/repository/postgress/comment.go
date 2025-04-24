package postgress

import (
	"1337b04rd/internal/domain/entity"
	"context"
	"database/sql"
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
	var comments []*entity.Comment
	return comments, nil
}

func (r *PostgresCommentRepository) DeleteComment(ctx context.Context, id int) error {
	return nil
}
