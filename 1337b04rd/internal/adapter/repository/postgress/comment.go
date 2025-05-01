package postgress

import (
	"a1337b04rd/internal/domain/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	query := `
		INSERT INTO comments (content, file_url, created_at, post_id, user_id, parent_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	// Если ParentID не задан (равен 0), передаем nil для записи NULL в базу
	var parentID interface{}
	if comment.ParentID != 0 {
		parentID = comment.ParentID
	} else {
		parentID = nil
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.FileURL,
		comment.CreatedAt,
		comment.PostID,
		comment.UserID,
		parentID,
	).Scan(&comment.ID)
	fmt.Println(parentID)
	if err != nil {
		return nil, fmt.Errorf("не удалось вставить комментарий: %w", err)
	}

	return comment, nil
}

func (r *PostgresCommentRepository) GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, user_id, content, created_at, parent_id, file_url
          FROM comments
         WHERE post_id = $1
         ORDER BY created_at ASC
    `, postID)
	if err != nil {
		log.Printf("Error executing query: %v", err) // Выводим ошибку при выполнении запроса
		return nil, err
	}
	defer rows.Close()

	var out []*entity.Comment
	for rows.Next() {
		var c entity.Comment
		var parentID sql.NullInt64 // используем sql.NullInt64 для работы с NULL значениями

		if err := rows.Scan(&c.ID, &c.UserID, &c.Content, &c.CreatedAt, &parentID, &c.FileURL); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Проверяем значение parentID и присваиваем 0, если NULL
		if parentID.Valid {
			c.ParentID = int(parentID.Int64) // преобразуем в int, если valid
		} else {
			c.ParentID = 0 // Если parentID NULL, присваиваем 0
		}
		c.FileURL = strings.Replace(c.FileURL, "minio", "localhost", 1)
		c.FileURL = strings.Replace(c.FileURL, "https", "http", 1)

		out = append(out, &c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	return out, nil
}

func (r *PostgresCommentRepository) DeleteComment(ctx context.Context, id int) error {
	return nil
}
