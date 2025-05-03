package servicework

import (
	"context"
	"fmt"
	"time"

	"a1337b04rd/internal/domain/entity"
	"a1337b04rd/internal/domain/port/repository"
)

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) *sessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *sessionService) CreateSession(ctx context.Context, session entity.Session) (entity.Session, error) {
	fmt.Println(">>> CreateSession service called")

	// Генерируем ID только если он не был установлен ранее
	if session.ID == "" {
		session.ID = generateSessionID()
	}

	// Устанавливаем время создания и истечения
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}

	if session.ExpiresAt.IsZero() {
		session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	}

	fmt.Printf(">>> Сохраняем сессию в БД: ID=%s, UserID=%d\n", session.ID, session.UserID)

	// Сохраняем в репозитории
	createdSession, err := s.sessionRepo.CreateSession(ctx, &session)
	if err != nil {
		fmt.Printf(">>> Ошибка сохранения сессии: %v\n", err)
		return entity.Session{}, err
	}

	fmt.Printf(">>> Сессия успешно создана: ID=%s\n", createdSession.ID)
	return *createdSession, nil
}

func (s *sessionService) GetSessionByID(ctx context.Context, id string) (entity.Session, error) {
	fmt.Printf(">>> GetSessionByID service called with ID: %s\n", id)

	session, err := s.sessionRepo.GetSessionByID(ctx, id)
	if err != nil {
		fmt.Printf(">>> Ошибка получения сессии: %v\n", err)
		return entity.Session{}, fmt.Errorf("failed to get session with ID %s: %v", id, err)
	}

	fmt.Printf(">>> Сессия найдена: UserID=%d, ExpiresAt=%v\n", session.UserID, session.ExpiresAt)
	return *session, nil
}

func (s *sessionService) DeleteSession(ctx context.Context, id string) error {
	err := s.sessionRepo.DeleteSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get session with ID %s: %v", id, err)
	}
	return nil
}

func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
