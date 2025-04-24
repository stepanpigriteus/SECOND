package handler

import (
	"1337b04rd/internal/domain/service"
	"fmt"
	"net/http"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> CreateSession handler called")
}

func (h *SessionHandler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> GetSessionById  handler called")
}

func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> GetSessionById  handler called")
}
