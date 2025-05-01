package handler

import (
	"1337b04rd/internal/domain/service"
	"1337b04rd/internal/api"
)

type AllHandlers struct {
	Post    *PostHandler
	User    *UserHandler
	Comment *CommentHandler
	Session *SessionHandler
}

func NewAllHandlers(services *service.AllServices) *AllHandlers {
	return &AllHandlers{
		Post:    NewPostHandler(services.Post, services.User, services.Comment, storage.FileStorage),
		User:    NewUserHandler(services.User),
		Comment: NewCommentHandler(services.Comment),
		Session: NewSessionHandler(services.Session),
	}
}

type contextKey string

const userContextKey contextKey = "user"
