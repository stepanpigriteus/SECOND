package servicework

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"a1337b04rd/internal/domain/entity"
)

// --- Мок для CommentRepository ---

type mockCommentRepo struct {
	createErr    error
	created      *entity.Comment
	getErr       error
	commentsByID []*entity.Comment
	deleteErr    error
}

func (m *mockCommentRepo) CreateComment(ctx context.Context, c *entity.Comment) (*entity.Comment, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	// копируем данные для внешней проверки
	*m.created = *c
	return c, nil
}

func (m *mockCommentRepo) GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.commentsByID, nil
}

func (m *mockCommentRepo) DeleteComment(ctx context.Context, id int) error {
	return m.deleteErr
}

// --- Мок для PostRepository ---

type mockPostRepo struct {
	updateErr error
	updatedID int
}

func (m *mockPostRepo) CreatePost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	return post, nil
}

func (m *mockPostRepo) GetPostByID(ctx context.Context, id int32) (*entity.Post, error) {
	return &entity.Post{}, nil
}

func (m *mockPostRepo) UpdatePost(ctx context.Context, postID int) error {
	m.updatedID = postID
	return m.updateErr
}

func (m *mockPostRepo) DeletePost(ctx context.Context, id int32) error {
	return nil
}

func (m *mockPostRepo) ListPosts(ctx context.Context) ([]entity.PostRequest, error) {
	return nil, nil
}

func (m *mockPostRepo) ListArchivedPosts(ctx context.Context) ([]entity.PostRequest, error) {
	return nil, nil
}

// --- Тесты commentService ---

func TestCreateComment_Success(t *testing.T) {
	// подготовка
	stored := &entity.Comment{}
	mc := &mockCommentRepo{created: stored}
	mp := &mockPostRepo{}
	svc := NewCommentService(mc, mp)

	input := entity.Comment{PostID: 7, Content: "Hello"}
	out, err := svc.CreateComment(context.Background(), input, (*multipart.FileHeader)(nil))
	if err != nil {
		t.Fatalf("CreateComment вернул ошибку: %v", err)
	}

	// CreatedAt установлен
	if out.CreatedAt.IsZero() {
		t.Error("ожидали, что CreatedAt установлен, получили zero")
	}
	// репозиторий получил тот же коммент
	if stored.PostID != input.PostID || stored.Content != input.Content {
		t.Errorf("репозиторий получил неверные данные: %+v", stored)
	}
	// UpdatePost вызван
	if mp.updatedID != input.PostID {
		t.Errorf("ожидали вызов UpdatePost с %d, получили %d", input.PostID, mp.updatedID)
	}
}

func TestCreateComment_CreateError(t *testing.T) {
	mc := &mockCommentRepo{createErr: errors.New("create fail")}
	mp := &mockPostRepo{}
	svc := NewCommentService(mc, mp)

	_, err := svc.CreateComment(context.Background(), entity.Comment{PostID: 1}, nil)
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if err.Error() != "create fail" {
		t.Errorf("ожидали 'create fail', получили %v", err)
	}
}

func TestCreateComment_UpdateError(t *testing.T) {
	stored := &entity.Comment{}
	mc := &mockCommentRepo{created: stored}
	mp := &mockPostRepo{updateErr: errors.New("upd fail")}
	svc := NewCommentService(mc, mp)

	_, err := svc.CreateComment(context.Background(), entity.Comment{PostID: 5}, nil)
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if !errors.Is(err, mp.updateErr) || !strings.Contains(err.Error(), "не удалось обновить время последнего комментария") {
		t.Errorf("ожидали обёртку 'не удалось обновить...', получили %v", err)
	}
}

func TestGetCommentByID_Success(t *testing.T) {
	now := time.Now()
	c1 := &entity.Comment{ID: 1, PostID: 3, Content: "A", CreatedAt: now}
	c2 := &entity.Comment{ID: 2, PostID: 3, Content: "B", CreatedAt: now}
	mc := &mockCommentRepo{commentsByID: []*entity.Comment{c1, c2}}
	svc := NewCommentService(mc, nil)

	out, err := svc.GetCommentByID(context.Background(), 3)
	if err != nil {
		t.Fatalf("GetCommentByID вернул ошибку: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("ожидали 2 комментария, получили %d", len(out))
	}
	if out[0].ID != 1 || out[1].ID != 2 {
		t.Errorf("неверные комментарии: %+v", out)
	}
}

func TestGetCommentByID_NotFound(t *testing.T) {
	mc := &mockCommentRepo{commentsByID: []*entity.Comment{}}
	svc := NewCommentService(mc, nil)

	_, err := svc.GetCommentByID(context.Background(), 10)
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	expected := fmt.Sprintf("comment not found for postID %d", 10)
	if err.Error() != expected {
		t.Errorf("ожидали '%s', получили '%v'", expected, err)
	}
}

func TestGetCommentByID_RepoError(t *testing.T) {
	mc := &mockCommentRepo{getErr: errors.New("db fail")}
	svc := NewCommentService(mc, nil)

	_, err := svc.GetCommentByID(context.Background(), 1)
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if err.Error() != "db fail" {
		t.Errorf("ожидали 'db fail', получили '%v'", err)
	}
}

func TestDeleteComment_SuccessAndError(t *testing.T) {
	mc := &mockCommentRepo{}
	svc := NewCommentService(mc, nil)

	// успешно
	if err := svc.DeleteComment(context.Background(), 11); err != nil {
		t.Errorf("DeleteComment вернул ошибку: %v", err)
	}

	// ошибка репозитория
	mc.deleteErr = errors.New("del fail")
	err := svc.DeleteComment(context.Background(), 11)
	if err == nil {
		t.Fatal("ожидали ошибку, получили nil")
	}
	if err.Error() != "del fail" {
		t.Errorf("ожидали 'del fail', получили '%v'", err)
	}
}
