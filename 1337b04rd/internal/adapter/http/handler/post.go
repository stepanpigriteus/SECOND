package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/service"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> CreatePost handler called")
	var post entity.Post
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	// Вызов слоя сервиса (с контекстом)
	createdPost, err := h.postService.CreatePost(ctx, post)
	if err != nil {
		http.Error(w, `{"error":"failed to create post"}`, http.StatusInternalServerError)
		return
	}

	// Возврат результата клиенту
	if err := json.NewEncoder(w).Encode(createdPost); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}
