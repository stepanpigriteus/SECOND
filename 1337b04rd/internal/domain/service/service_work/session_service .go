package servicework

import (
	"context"
	"fmt"
	"time"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port/repository"
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
	session.ID = generateSessionID()
	session.CreatedAt = time.Now()
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	createdSession, err := s.sessionRepo.CreateSession(ctx, &session)
	if err != nil {
		return entity.Session{}, err
	}
	return *createdSession, nil
}

func (s *sessionService) GetSessionByID(ctx context.Context, id string) (entity.Session, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, id)
	if err != nil {
		return entity.Session{}, fmt.Errorf("failed to get session with ID %s: %v", id, err)
	}
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
