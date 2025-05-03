package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"a1337b04rd/internal/adapter/http/middleware"
	"a1337b04rd/internal/domain/entity"
)

// --- Моки ---
type mockSessionService struct{}

func (m *mockSessionService) GetSessionByID(ctx context.Context, id string) (entity.Session, error) {
	return entity.Session{}, nil
}

func (m *mockSessionService) CreateSession(ctx context.Context, s entity.Session) (entity.Session, error) {
	return entity.Session{
		ID:        "test-session-id",
		UserID:    s.UserID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}, nil
}

// Добавляем метод DeleteSession
func (m *mockSessionService) DeleteSession(ctx context.Context, id string) error {
	return nil
}

type mockUserService struct{}

func (m *mockUserService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	return entity.User{ID: id, CharacterName: "Test", AvatarURL: "https://example.com"}, nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, u entity.User) (entity.User, error) {
	return u, nil
}

func (m *mockUserService) CreateUser(ctx context.Context, u entity.User) (entity.User, error) {
	return entity.User{ID: 123, CharacterName: "Test", AvatarURL: "https://example.com"}, nil
}

// Добавляем метод DeleteUser
func (m *mockUserService) DeleteUser(ctx context.Context, id int) error {
	return nil
}

// --- Тест ---

func TestSessionMiddleware_NewUserFlow(t *testing.T) {
	var nextCalled bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true

		// Проверяем, что пользователь есть в контексте
		user, ok := r.Context().Value("user").(entity.User)
		if !ok {
			t.Errorf("user not in context or wrong type")
			return
		}
		if user.ID != 123 {
			t.Errorf("expected user ID 123, got %d", user.ID)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
	rr := httptest.NewRecorder()

	middlewareToTest := middleware.SessionMiddleware(
		&mockSessionService{},
		&mockUserService{},
	)

	middlewareToTest(handler).ServeHTTP(rr, req)

	if !nextCalled {
		t.Errorf("next handler was not called")
	}

	// Проверка установки куки
	cookies := rr.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "session_id" && c.Value == "test-session-id" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("session_id cookie not set or incorrect")
	}
}
