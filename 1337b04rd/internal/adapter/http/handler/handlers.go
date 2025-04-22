package handler

import (
	"1337b04rd/internal/domain/service"
)

type AllHandlers struct {
	Post *PostHandler
	// 	User    *UserHandler
	// 	Comment *CommentHandler
	// 	Session *SessionHandler
}

func NewAllHandlers(services *service.AllServices) *AllHandlers {
	return &AllHandlers{
		Post: NewPostHandler(&services.Post),
		// User:    NewUserHandler(services.User),
		// Comment: NewCommentHandler(services.Comment),
		// Session: NewSessionHandler(services.Session),
	}
}
