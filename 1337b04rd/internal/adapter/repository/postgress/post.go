package postgress

import (
	"1337b04rd/internal/domain/entity"
	"context"
	"database/sql"
	"fmt"
)

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

func (r *PostgresPostRepository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	var post entity.Post

	///заглушка
	query := `SELECT id, title, content, image_url, user_id, created_at, updated_at, deleted_at, last_comment_at FROM posts WHERE id = $1`

	// Выполняем запрос
	row := r.db.QueryRowContext(ctx, query, id)

	// Заполняем структуру post данными из базы
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt, &post.LastCommentAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return &post, nil
}

func (r *PostgresPostRepository) UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	return post, nil
}

func (r *PostgresPostRepository) DeletePost(ctx context.Context, id int) error {
	return nil
}
