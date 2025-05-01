package handler

import (
	"a1337b04rd/internal/domain/entity"
	"a1337b04rd/internal/domain/service"
	"encoding/json"
	"net/http"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var session entity.Session
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, "неправильный формат тела запроса", http.StatusBadRequest)
		return
	}

	createdSession, err := h.sessionService.CreateSession(r.Context(), session)
	if err != nil {
		http.Error(w, "ошибка при создании сессии: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSession)
}

func (h *SessionHandler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "не передан ID сессии", http.StatusBadRequest)
		return
	}

	session, err := h.sessionService.GetSessionByID(r.Context(), id)
	if err != nil {
		http.Error(w, "ошибка при получении сессии: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(session)
}

func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "не передан ID сессии", http.StatusBadRequest)
		return
	}

	err := h.sessionService.DeleteSession(r.Context(), id)
	if err != nil {
		http.Error(w, "ошибка при удалении сессии: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
