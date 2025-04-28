package postgress

import (
	"context"
	"database/sql"
	"fmt"

	"1337b04rd/internal/domain/entity"
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

func (r *PostgresPostRepository) GetPostByID(ctx context.Context, id int32) (*entity.Post, error) {
	var post entity.Post

	///заглушка
	query := `SELECT id, title, content, image_url, user_id, created_at, updated_at, COALESCE(deleted_at, '1970-01-01 00:00:00') AS deleted_at, last_comment_at FROM posts WHERE id = $1;`
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt, &post.LastCommentAt)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows returned from query")
			return nil, fmt.Errorf("post not found")
		}
		fmt.Println("Error during scan:", err)
		return nil, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return &post, nil
}

func (r *PostgresPostRepository) UpdatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	return post, nil
}

func (r *PostgresPostRepository) DeletePost(ctx context.Context, id int32) error {
	return nil
}

func (r *PostgresPostRepository) ListPosts(ctx context.Context) ([]entity.PostRequest, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, title, image_url
        FROM posts
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.PostRequest
	for rows.Next() {
		var post entity.PostRequest
		if err := rows.Scan(&post.ID, &post.Title, &post.ImageURL); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
