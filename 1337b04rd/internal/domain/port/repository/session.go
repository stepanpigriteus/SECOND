package repository

import "1337b04rd/internal/domain/entity"

type SessionRepository interface {
	CreateSession(session *entity.Session) (*entity.Session, error)
	GetSessionByID(id string) (*entity.Session, error)
	DeleteSession(id string) error
}
