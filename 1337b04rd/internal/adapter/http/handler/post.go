package handler

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/service"
	externalfunc "1337b04rd/pkg/external_func"
	"encoding/json"
	"fmt"
	"net/http"
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
		http.Error(w, `{"error":"invalid input!"}`, http.StatusBadRequest)
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

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> GetPostByID handler called")
	if h.postService == nil {
		http.Error(w, `{"error":"postService is not initialized"}`, http.StatusInternalServerError)
		return
	}

	id, err := externalfunc.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id32 := int32(id)
	post, err := h.postService.GetPostByID(ctx, id32)
	if err != nil {
		http.Error(w, `{"error":"post not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdatePostPostByID(w http.ResponseWriter, r *http.Request) {
}

func (h *PostHandler) DeletePostPostPostByID(w http.ResponseWriter, r *http.Request) {
}
