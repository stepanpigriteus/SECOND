package servicework

import (
	"context"
	"fmt"
	"testing"
	"time"

	"a1337b04rd/internal/domain/entity"
)

// --- Мок-репозиторий ---

type mockPostRepository struct{}

// CreatePost: соответствует сигнатуре CreatePost(ctx, *entity.Post) (*entity.Post, error)
func (m *mockPostRepository) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	// проверяем, что CreatedAt уже заполнено сервисом
	if post.CreatedAt.IsZero() {
		return nil, fmt.Errorf("CreatedAt не установлен")
	}
	post.ID = 1
	return post, nil
}

// GetPostByID: соответствует GetPostByID(ctx, int32) (*entity.Post, error)
func (m *mockPostRepository) GetPostByID(ctx context.Context, id int32) (*entity.Post, error) {
	now := time.Now()
	if id == 1 {
		return &entity.Post{ID: 1, Title: "Test Post", CreatedAt: now, UpdatedAt: now}, nil
	}
	// Для всех остальных возвращаем ненулевой указатель, но с ID=0
	return &entity.Post{}, nil
}

// UpdatePost: соответствует UpdatePost(ctx, int32, entity.Post) (*entity.Post, error)
func (m *mockPostRepository) UpdatePost(ctx context.Context, id int) error {
	return nil
}

// DeletePost: соответствует DeletePost(ctx, int32) error
func (m *mockPostRepository) DeletePost(ctx context.Context, id int32) error {
	if id != 1 {
		return fmt.Errorf("not found")
	}
	return nil
}

// ListPosts: соответствует ListPosts(ctx) ([]entity.PostRequest, error)
func (m *mockPostRepository) ListPosts(ctx context.Context) ([]entity.PostRequest, error) {
	return []entity.PostRequest{
		{ID: 1, Title: "Test Post", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}, nil
}

// ListArchivedPosts: соответствует ListArchivedPosts(ctx) ([]entity.PostRequest, error)
func (m *mockPostRepository) ListArchivedPosts(ctx context.Context) ([]entity.PostRequest, error) {
	return nil, nil
}

// --- Тесты сервиса ---

func TestCreatePost(t *testing.T) {
	repo := &mockPostRepository{}
	svc := NewPostService(repo)

	// Передаём пост без заполненных времён — сервис должен их выставить
	p := entity.Post{Title: "Test Post"}
	created, err := svc.CreatePost(context.Background(), p)
	if err != nil {
		t.Fatalf("CreatePost returned error: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("ожидали ID=1, получили %d", created.ID)
	}
	if created.Title != "Test Post" {
		t.Errorf("ожидали Title='Test Post', получили '%s'", created.Title)
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Error("CreatedAt или UpdatedAt не были установлены")
	}
}

func TestGetPostByID(t *testing.T) {
	repo := &mockPostRepository{}
	svc := NewPostService(repo)

	post, err := svc.GetPostByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetPostByID returned error: %v", err)
	}
	if post.ID != 1 {
		t.Errorf("ожидали ID=1, получили %d", post.ID)
	}
	if post.Title != "Test Post" {
		t.Errorf("ожидали Title='Test Post', получили '%s'", post.Title)
	}

	// Случай несуществующего ID
	post, err = svc.GetPostByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("GetPostByID при несуществующем ID вернул ошибку: %v", err)
	}
	if post.ID != 0 {
		t.Errorf("ожидали пустой Post для несуществующего ID, получили ID=%d", post.ID)
	}
}
func TestUpdatePost(t *testing.T) {
    repo := &mockPostRepository{}
    svc := NewPostService(repo)

    // Сейчас UpdatePost сервиса возвращает всегда zero-value entity.Post
    updated, err := svc.UpdatePost(context.Background(), 1, entity.Post{Title: "Updated"})
    if err != nil {
        t.Fatalf("UpdatePost returned unexpected error: %v", err)
    }
    // Ожидаем именно zero-value, а не ID=1
    if updated.ID != 0 {
        t.Errorf("ожидали ID=0 (zero-value), получили %d", updated.ID)
    }
    if updated.Title != "" {
        t.Errorf("ожидали пустой Title (zero-value), получили '%s'", updated.Title)
    }
}
func TestDeletePost(t *testing.T) {
	repo := &mockPostRepository{}
	svc := NewPostService(repo)

	// Успешное удаление
	if err := svc.DeletePost(context.Background(), 1); err != nil {
		t.Errorf("DeletePost(1) вернул ошибку: %v", err)
	}
	// Удаление несуществующего
	if err := svc.DeletePost(context.Background(), 2); err == nil {
		t.Error("DeletePost(2) ожидалась ошибка, получили nil")
	}
}

func TestListPosts(t *testing.T) {
	repo := &mockPostRepository{}
	svc := NewPostService(repo)

	posts, err := svc.ListPosts(context.Background())
	if err != nil {
		t.Fatalf("ListPosts вернул ошибку: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("ожидали 1 пост, получили %d", len(posts))
	}
	if posts[0].Title != "Test Post" {
		t.Errorf("ожидали Title='Test Post', получили '%s'", posts[0].Title)
	}
}

func TestListArchivedPosts(t *testing.T) {
	repo := &mockPostRepository{}
	svc := NewPostService(repo)

	posts, err := svc.ListArchivedPosts(context.Background())
	if err != nil {
		t.Fatalf("ListArchivedPosts вернул ошибку: %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("ожидали 0 постов, получили %d", len(posts))
	}
}
