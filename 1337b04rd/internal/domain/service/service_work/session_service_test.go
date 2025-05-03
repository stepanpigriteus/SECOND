package servicework

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"a1337b04rd/internal/domain/entity"
)

// --- Мок-репозиторий ---

type mockSessionRepo struct {
	// для симуляции ошибок
	createErr   error
	getErr      error
	deleteErr   error
	stored      *entity.Session
	returnEmpty bool
}

func (m *mockSessionRepo) CreateSession(ctx context.Context, s *entity.Session) (*entity.Session, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	// симулируем сохранение
	*m.stored = *s
	return s, nil
}

func (m *mockSessionRepo) GetSessionByID(ctx context.Context, id string) (*entity.Session, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.returnEmpty {
		// вернуть пустую, но не nil
		zero := entity.Session{}
		return &zero, nil
	}
	return m.stored, nil
}

func (m *mockSessionRepo) DeleteSession(ctx context.Context, id string) error {
	return m.deleteErr
}

// --- Тесты ---

func TestCreateSession_Defaults(t *testing.T) {
	// Подготовка моков
	stored := &entity.Session{}
	mock := &mockSessionRepo{stored: stored}
	svc := NewSessionService(mock)

	// Вызываем CreateSession с пустым полем ID, CreatedAt, ExpiresAt
	input := entity.Session{UserID: 42}
	created, err := svc.CreateSession(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateSession returned unexpected error: %v", err)
	}

	// Проверяем, что ID не пустой
	if created.ID == "" {
		t.Error("expected non-empty session ID, got empty")
	}
	// Проверяем, что CreatedAt и ExpiresAt установлены
	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set, got zero")
	}
	if created.ExpiresAt.IsZero() {
		t.Error("expected ExpiresAt to be set, got zero")
	}
	if !created.ExpiresAt.After(created.CreatedAt) {
		t.Errorf("expected ExpiresAt > CreatedAt, got CreatedAt=%v, ExpiresAt=%v", created.CreatedAt, created.ExpiresAt)
	}
	// Проверяем, что репозиторий получил ту же структуру
	if stored.ID != created.ID || stored.UserID != created.UserID {
		t.Errorf("repository stored data mismatch: got %+v", stored)
	}
}

func TestCreateSession_KeepExistingID(t *testing.T) {
	stored := &entity.Session{}
	mock := &mockSessionRepo{stored: stored}
	svc := NewSessionService(mock)

	// Задаём своё ID и CreatedAt/ExpiresAt
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	input := entity.Session{
		ID:        "custom-id",
		UserID:    7,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}
	created, err := svc.CreateSession(context.Background(), input)
	if err != nil {
		t.Fatalf("CreateSession returned unexpected error: %v", err)
	}

	// ID и времена не должны измениться
	if created.ID != "custom-id" {
		t.Errorf("expected ID 'custom-id', got '%s'", created.ID)
	}
	if !created.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, created.CreatedAt)
	}
	if !created.ExpiresAt.Equal(now.Add(time.Hour)) {
		t.Errorf("expected ExpiresAt %v, got %v", now.Add(time.Hour), created.ExpiresAt)
	}
}

func TestCreateSession_Error(t *testing.T) {
	mock := &mockSessionRepo{
		stored:    &entity.Session{},
		createErr: errors.New("db error"),
	}
	svc := NewSessionService(mock)

	_, err := svc.CreateSession(context.Background(), entity.Session{UserID: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "db error") {
		t.Errorf("expected error to contain 'db error', got %v", err)
	}
}

func TestGetSessionByID_Success(t *testing.T) {
	now := time.Now()
	stored := &entity.Session{ID: "sid", UserID: 5, CreatedAt: now, ExpiresAt: now.Add(time.Hour)}
	mock := &mockSessionRepo{stored: stored}
	svc := NewSessionService(mock)

	got, err := svc.GetSessionByID(context.Background(), "sid")
	if err != nil {
		t.Fatalf("GetSessionByID returned error: %v", err)
	}
	if got.ID != "sid" || got.UserID != 5 {
		t.Errorf("unexpected session returned: %+v", got)
	}
}

func TestGetSessionByID_NotFound(t *testing.T) {
	stored := &entity.Session{ID: "other", UserID: 0}
	mock := &mockSessionRepo{stored: stored, returnEmpty: true}
	svc := NewSessionService(mock)

	got, err := svc.GetSessionByID(context.Background(), "sid")
	if err != nil {
		t.Fatalf("GetSessionByID returned error: %v", err)
	}
	// returnEmpty даёт zero-value
	if got.ID != "" && got.UserID != 0 {
		t.Errorf("expected zero-value session, got %+v", got)
	}
}

func TestGetSessionByID_Error(t *testing.T) {
	mock := &mockSessionRepo{getErr: errors.New("not found")}
	svc := NewSessionService(mock)

	_, err := svc.GetSessionByID(context.Background(), "sid")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get session with ID sid") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDeleteSession_Success(t *testing.T) {
	mock := &mockSessionRepo{}
	svc := NewSessionService(mock)

	if err := svc.DeleteSession(context.Background(), "sid"); err != nil {
		t.Errorf("DeleteSession returned error: %v", err)
	}
}

func TestDeleteSession_Error(t *testing.T) {
	mock := &mockSessionRepo{deleteErr: errors.New("del fail")}
	svc := NewSessionService(mock)

	err := svc.DeleteSession(context.Background(), "sid")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get session with ID sid") {
		t.Errorf("unexpected error message: %v", err)
	}
}
