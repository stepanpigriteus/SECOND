package handler

import (
	"a1337b04rd/internal/domain/entity"
	"a1337b04rd/internal/domain/service"
	externalfunc "a1337b04rd/pkg/external_func"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> CreateComment handler called")

	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	// Ограничить размер тела (например, 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"invalid multipart form"}`, http.StatusBadRequest)
		return
	}

	text := r.FormValue("comment")
	file, fileHeader, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, `{"error":"failed to read file"}`, http.StatusBadRequest)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	postIDStr := r.FormValue("post_id")

	fmt.Println(">>> post ID", postIDStr)
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post_id", http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user")
	parID := r.FormValue("parent_id")
	var parentID int
	if parID != "" {
		parentID, err = strconv.Atoi(parID)
	}

	if err != nil {

		http.Error(w, "invalid parent_id", http.StatusBadRequest)
		return
	}

	userID := user.(*entity.User).ID
	comment := entity.Comment{
		Content:  text,
		PostID:   postID,
		UserID:   userID,
		ParentID: parentID,
		// Другие поля можно добавить по необходимости
	}

	// Вызов сервиса
	createdComment, err := h.commentService.CreateComment(ctx, comment, fileHeader)
	if err != nil {
		http.Error(w, `{"error":"failed to create comment"}`, http.StatusInternalServerError)
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
