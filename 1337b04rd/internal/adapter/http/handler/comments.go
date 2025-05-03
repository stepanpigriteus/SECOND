package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"a1337b04rd/internal/domain/entity"
	"a1337b04rd/internal/domain/port"
	apiport "a1337b04rd/internal/domain/port/api"
	"a1337b04rd/internal/domain/service"
	externalfunc "a1337b04rd/pkg/external_func"
)

type CommentHandler struct {
	commentService service.CommentService
	postService    service.PostService
	fileStorage    apiport.FileStorage
	logger         port.Logger
}

func NewCommentHandler(commentService service.CommentService, fs apiport.FileStorage, postService service.PostService, logger port.Logger) *CommentHandler {
	return &CommentHandler{commentService: commentService, fileStorage: fs, postService: postService, logger: logger}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> CreateComment handler called")
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	// Ограничить размер тела (например, 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Warn("invalid multipart form:", err)
		http.Error(w, `{"error":"invalid multipart form or file too large"}`, http.StatusBadRequest)
		return
	}

	text := r.FormValue("comment")
	postIDStr := r.FormValue("post_id")
	parID := r.FormValue("parent_id")

	file, handler, err := r.FormFile("file")
	var urlFile string
	if err != nil && err != http.ErrMissingFile {
		h.logger.Warn("error reading uploaded file:", err)
		http.Error(w, `{"error":"invalid uploaded file"}`, http.StatusBadRequest)
		return
	}

	if err == nil && handler != nil {
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			h.logger.Warn("failed to read file bytes:", err)
			http.Error(w, `{"error":"failed to read file"}`, http.StatusInternalServerError)
			return
		}

		fileReader := strings.NewReader(string(fileBytes))
		objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)

		urlFile, err = h.fileStorage.UploadFile(ctx, "comments", objectName, fileReader, int64(len(fileBytes)), handler.Header.Get("Content-Type"))
		if err != nil {
			h.logger.Error("file upload failed:", err)
			http.Error(w, `{"error":"failed to upload file"}`, http.StatusInternalServerError)
			return
		}
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid post_id"}`, http.StatusBadRequest)
		return
	}

	var parentID int
	if parID != "" {
		parentID, err = strconv.Atoi(parID)
		if err != nil {
			http.Error(w, `{"error":"invalid parent_id"}`, http.StatusBadRequest)
			return
		}
	}

	user := ctx.Value("user")
	userID := user.(*entity.User).ID

	comment := entity.Comment{
		Content:  text,
		PostID:   postID,
		UserID:   userID,
		ParentID: parentID,
		FileURL:  urlFile,
	}

	createdComment, err := h.commentService.CreateComment(ctx, comment, handler)
	if err != nil {
		http.Error(w, `{"error":"failed to create comment"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(createdComment); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *CommentHandler) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> GetCommentByID handler called")
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
	if err != nil {
		http.Error(w, `{"error":"error comment extraction"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		h.logger.Error("error: failed to encode response")
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
}
