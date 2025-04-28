package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port"
	"1337b04rd/internal/domain/service"
)

type PostHandler struct {
	postService    service.PostService
	userService    service.UserService
	commentService service.CommentService
	logger         port.Logger
}

func NewPostHandler(ps service.PostService, us service.UserService, cs service.CommentService) *PostHandler {
	return &PostHandler{
		postService:    ps,
		userService:    us,
		commentService: cs,
	}
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
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	id64, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	postID := int32(id64)

	ctx := r.Context()

	post, err := h.postService.GetPostByID(ctx, postID)
	if err != nil {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	user, err := h.userService.GetUserByID(ctx, post.UserID)
	if err != nil {
		http.Error(w, "author not found", http.StatusInternalServerError)
		return
	}

	comments, err := h.commentService.GetCommentByID(r.Context(), int(postID))
	if err != nil {
		http.Error(w, "failed to load comments", http.StatusInternalServerError)
		return
	}

	cmVM := make([]CommentVM, 0, len(comments))
	for _, c := range comments {
		cu, err := h.userService.GetUserByID(ctx, c.UserID)
		if err != nil {
			continue
		}
		cmVM = append(cmVM, CommentVM{
			AvatarURL: cu.AvatarURL,
			UserName:  cu.CustomName,
			CreatedAt: c.CreatedAt.Format(time.RFC1123),
			CommentID: int32(c.ID),
			Content:   c.Content,
		})
	}

	// 6) Собираем итоговый VM
	vm := PostVM{
		UserAvatar: user.AvatarURL,
		UserName:   user.CustomName,
		CreatedAt:  post.CreatedAt.Format(time.RFC1123),
		PostID:     post.ID,
		Title:      post.Title,
		Content:    post.Content,
		Comments:   cmVM,
		Image:      post.ImageURL,
	}

	// 7) Рендерим шаблон
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := RenderTemplate(w, "post.html", vm); err != nil {
		http.Error(w, "Error template rendering: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdatePostPostByID(w http.ResponseWriter, r *http.Request) {
}

func (h *PostHandler) DeletePostPostPostByID(w http.ResponseWriter, r *http.Request) {
}

func (h *PostHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := h.postService.ListPosts(ctx)
	if err != nil {
		h.logger.Error("Ошибка при получении постов", http.StatusInternalServerError)
		http.Error(w, `{"error":"error receiving posts"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем данные в шаблон
	w.Header().Set("Content-Type", "text/html")
	err = RenderTemplate(w, "catalog.html", posts)
	if err != nil {
		http.Error(w, "Error template rendering", http.StatusInternalServerError)
		return
	}
}
