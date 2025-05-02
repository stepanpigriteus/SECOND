package handler

import (
	"a1337b04rd/internal/domain/service"
)

type AllHandlers struct {
	Post    *PostHandler
	User    *UserHandler
	Comment *CommentHandler
	Session *SessionHandler
}

func NewAllHandlers(services *service.AllServices) *AllHandlers {
	return &AllHandlers{
		Post:    NewPostHandler(services.Post, services.User, services.Comment, services.Storage),
		User:    NewUserHandler(services.User),
		Comment: NewCommentHandler(services.Comment, services.Storage, services.Post),
		Session: NewSessionHandler(services.Session),
	}
}

type contextKey string

const userContextKey contextKey = "user"
