package postgress

import (
	"context"
	"database/sql"

	"1337b04rd/internal/domain/entity"
)

// type PostRepository interface {
//     CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
//     GetPostByID(ctx context.Context, id int) (*entity.Post, error)
//     UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error)
//     DeletePost(ctx context.Context, id int) error
// }

type PostgresPostRepository struct {
	db *sql.DB
}

func NewPostgresPostRepository(db *sql.DB) *PostgresPostRepository {
	return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	// запрос на создание поста
	return post, nil
}

func (r *PostgresPostRepository) GetPostByID(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	return post, nil
}

func (r *PostgresPostRepository) UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	return post, nil
}

func (r *PostgresPostRepository) DeletePost(ctx context.Context, id int) error {
	return nil
}
