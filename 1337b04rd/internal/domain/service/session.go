package service

import (
	"context"

	"a1337b04rd/internal/domain/entity"
)

type SessionService interface {
	CreateSession(ctx context.Context, session entity.Session) (entity.Session, error)
	GetSessionByID(ctx context.Context, id string) (entity.Session, error)
	DeleteSession(ctx context.Context, id string) error
}


