package handler

import (
	"a1337b04rd/internal/domain/port"
	"a1337b04rd/internal/domain/service"
)

type AllHandlers struct {
	Post    *PostHandler
	User    *UserHandler
	Comment *CommentHandler
	Session *SessionHandler
}

func NewAllHandlers(services *service.AllServices, logger port.Logger) *AllHandlers {
	return &AllHandlers{
		Post:    NewPostHandler(services.Post, services.User, services.Comment, services.Storage, logger),
		User:    NewUserHandler(services.User),
		Comment: NewCommentHandler(services.Comment, services.Storage, services.Post, logger),
		Session: NewSessionHandler(services.Session),
	}
}
