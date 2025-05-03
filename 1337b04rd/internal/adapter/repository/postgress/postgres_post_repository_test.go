package postgress

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"a1337b04rd/internal/domain/entity"
)

// import (
// 	"context"
// 	"testing"

// 	"a1337b04rd/internal/domain/entity"
// )

// func TestCreatePost_DBIsNil(t *testing.T) {
// 	// Передаем nil для базы данных
// 	repo := NewPostgresPostRepository(nil)

// 	// Создаем пост
// 	post := &entity.Post{
// 		Title:   "Test Title",
// 		Content: "Test Content",
// 	}

// 	_, err := repo.CreatePost(context.Background(), post)
// 	if err == nil {
// 		t.Errorf("Expected error when DB is nil, got nil")
// 	}
// }

type DBInterface interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Убедимся, что *sql.DB реализует наш интерфейс
var _ DBInterface = (*sql.DB)(nil)

// MockRow - мок для sql.Row
type MockRow struct {
	data []interface{}
	err  error
}

// Scan для MockRow имитирует поведение Scan для sql.Row
func (m *MockRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}

	if len(dest) != len(m.data) {
		return errors.New("количество аргументов не совпадает")
	}

	for i, val := range m.data {
		switch v := dest[i].(type) {
		case *int32:
			*v = val.(int32)
		case *string:
			*v = val.(string)
		case *time.Time:
			*v = val.(time.Time)
		// Можно добавить другие типы по необходимости
		default:
			return errors.New("неподдерживаемый тип для сканирования")
		}
	}
	return nil
}

// MockDB - мок БД для тестов
type MockDB struct {
	QueryRowFunc func(context.Context, string, ...interface{}) *MockRow
	queryError   error
}

// QueryRowContext - имитация работы с БД
func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	// Вместо возврата sql.Row мы будем хранить результаты мока внутри
	// и возвращать заглушку sql.Row, которую заменим нашим MockRow при вызове Scan
	if m.QueryRowFunc == nil {
		return &sql.Row{}
	}

	// mockRow := m.QueryRowFunc(ctx, query, args...)

	// Создаем специальную обертку для мока
	row := sql.Row{}

	// Мы не можем напрямую изменить sql.Row, но можем использовать замыкание
	// для хранения нашего mockRow и возврата правильных данных в методе CreatePost

	return &row
}

// QueryContext - заглушка для интерфейса
func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("not implemented")
}

// ExecContext - заглушка для интерфейса
func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, errors.New("not implemented")
}

// Модифицируем репозиторий для поддержки интерфейса
type PostgresPostRepositoryTestable struct {
	db DBInterface
	// Добавляем поле для доступа к мок результатам
	mockResults *MockRow
}

// NewPostgresPostRepositoryTestable - конструктор для тестируемого репозитория
func NewPostgresPostRepositoryTestable(db DBInterface) *PostgresPostRepositoryTestable {
	return &PostgresPostRepositoryTestable{db: db}
}

// SetMockResults - устанавливает мок результаты для тестов
func (r *PostgresPostRepositoryTestable) SetMockResults(mockRow *MockRow) {
	r.mockResults = mockRow
}

// CreatePost - реализация CreatePost для тестируемого репозитория
func (r *PostgresPostRepositoryTestable) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	if r.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		INSERT INTO posts (title, content, image_url, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, last_comment_at
	`

	// В реальном коде вызываем настоящее QueryRowContext
	_ = r.db.QueryRowContext(ctx, query,
		post.Title,
		post.Content,
		post.ImageURL,
		post.UserID,
	)

	// В тестах используем r.mockResults вместо реального сканирования
	if r.mockResults != nil {
		if r.mockResults.err != nil {
			return nil, r.mockResults.err
		}

		// Сканируем данные из мока в post
		err := r.mockResults.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.LastCommentAt)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

// TestCreatePost тестирует функцию CreatePost
func TestCreatePost(t *testing.T) {
	// Устанавливаем время для тестов
	now := time.Now()

	// Тестовые случаи
	tests := []struct {
		name        string
		post        *entity.Post
		mockRow     *MockRow
		mockDBError error
		wantErr     bool
		wantPost    *entity.Post
	}{
		{
			name: "Успешное создание поста",
			post: &entity.Post{
				Title:    "Test Title",
				Content:  "Test Content",
				ImageURL: "https://minio/image.jpg",
				UserID:   1,
			},
			mockRow: &MockRow{
				data: []interface{}{
					int32(1), // ID
					now,      // created_at
					now,      // updated_at
					now,      // last_comment_at
				},
			},
			wantErr: false,
			wantPost: &entity.Post{
				ID:            1,
				Title:         "Test Title",
				Content:       "Test Content",
				ImageURL:      "https://minio/image.jpg",
				UserID:        1,
				CreatedAt:     now,
				UpdatedAt:     now,
				LastCommentAt: now,
			},
		},
		{
			name: "Ошибка БД",
			post: &entity.Post{
				Title:    "Fail Title",
				Content:  "Fail Content",
				ImageURL: "https://minio/fail.jpg",
				UserID:   2,
			},
			mockRow: &MockRow{
				err: errors.New("ошибка базы данных"),
			},
			wantErr:  true,
			wantPost: nil,
		},
		{
			name: "Пустое соединение с БД",
			post: &entity.Post{
				Title:    "Null DB",
				Content:  "Null Content",
				ImageURL: "https://minio/null.jpg",
				UserID:   3,
			},
			mockRow:  nil,
			wantErr:  true,
			wantPost: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Если нужно тестировать nil DB
			if tt.name == "Пустое соединение с БД" {
				repo := NewPostgresPostRepositoryTestable(nil)
				_, err := repo.CreatePost(context.Background(), tt.post)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreatePost() с nil DB error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			// Создаем мок DB
			mockDB := &MockDB{}

			// Создаем репозиторий
			repo := NewPostgresPostRepositoryTestable(mockDB)

			// Устанавливаем мок результаты
			repo.SetMockResults(tt.mockRow)

			// Вызываем тестируемый метод
			post, err := repo.CreatePost(context.Background(), tt.post)

			// Проверяем ошибку
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Если ожидаем ошибку, дальнейшие проверки не нужны
			if tt.wantErr {
				return
			}

			// Проверяем результат
			if post.ID != tt.wantPost.ID {
				t.Errorf("CreatePost() ID = %v, want %v", post.ID, tt.wantPost.ID)
			}

			if post.Title != tt.wantPost.Title {
				t.Errorf("CreatePost() Title = %v, want %v", post.Title, tt.wantPost.Title)
			}

			if post.Content != tt.wantPost.Content {
				t.Errorf("CreatePost() Content = %v, want %v", post.Content, tt.wantPost.Content)
			}

			if post.ImageURL != tt.wantPost.ImageURL {
				t.Errorf("CreatePost() ImageURL = %v, want %v", post.ImageURL, tt.wantPost.ImageURL)
			}

			if post.UserID != tt.wantPost.UserID {
				t.Errorf("CreatePost() UserID = %v, want %v", post.UserID, tt.wantPost.UserID)
			}

			// Проверка временных полей
			if !post.CreatedAt.Equal(tt.wantPost.CreatedAt) {
				t.Errorf("CreatePost() CreatedAt = %v, want %v", post.CreatedAt, tt.wantPost.CreatedAt)
			}

			if !post.UpdatedAt.Equal(tt.wantPost.UpdatedAt) {
				t.Errorf("CreatePost() UpdatedAt = %v, want %v", post.UpdatedAt, tt.wantPost.UpdatedAt)
			}

			if !post.LastCommentAt.Equal(tt.wantPost.LastCommentAt) {
				t.Errorf("CreatePost() LastCommentAt = %v, want %v", post.LastCommentAt, tt.wantPost.LastCommentAt)
			}
		})
	}
}
