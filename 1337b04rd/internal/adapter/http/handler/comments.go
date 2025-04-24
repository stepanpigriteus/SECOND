package handler

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/service"
	externalfunc "1337b04rd/pkg/external_func"
	"encoding/json"
	"fmt"
	"net/http"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> CreateComment handler called")
	var comment entity.Comment
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, `{"error":"invalid input!"}`, http.StatusBadRequest)
		return
	}

	createdComment, err := h.commentService.CreateComment(ctx, comment)
	if err != nil {
		http.Error(w, `{"error":"failed to create post"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(createdComment); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *CommentHandler) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> GetCommentByID handler called")
	if h.commentService == nil {
		http.Error(w, `{"error":"commentService is not initialized"}`, http.StatusInternalServerError)
		return
	}
	id, err := externalfunc.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"comment not found"}`, http.StatusNotFound)
		return
	}

	ctx := r.Context()
	comment, err := h.commentService.GetCommentByID(ctx, id)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {

}
