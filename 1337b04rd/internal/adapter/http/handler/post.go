package handler

import (
	"encoding/json"
	"net/http"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/service"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	var post entity.Post
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
