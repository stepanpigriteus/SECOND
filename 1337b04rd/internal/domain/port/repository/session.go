package repository

import (
	"1337b04rd/internal/domain/entity"
	"context"
)

type SessionRepository interface {
	CreateSession(context.Context, *entity.Session) (*entity.Session, error)
	GetSessionByID(ctx context.Context, id string) (*entity.Session, error)
	DeleteSession(ctx context.Context, id string) error
}
