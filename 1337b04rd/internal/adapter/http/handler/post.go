package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"a1337b04rd/internal/domain/entity"
	"a1337b04rd/internal/domain/port"
	apiport "a1337b04rd/internal/domain/port/api"
	"a1337b04rd/internal/domain/service"
)

type PostHandler struct {
	postService    service.PostService
	userService    service.UserService
	commentService service.CommentService
	logger         port.Logger
	fileStorage    apiport.FileStorage
}

func NewPostHandler(ps service.PostService, us service.UserService, cs service.CommentService, fs apiport.FileStorage, logger port.Logger) *PostHandler {
	return &PostHandler{
		postService:    ps,
		userService:    us,
		commentService: cs,
		fileStorage:    fs,
		logger:         logger,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>> CreatePost handler called")
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	// Получение пользователя из контекста
	userCtx := ctx.Value("user")
	if userCtx == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	user, ok := userCtx.(*entity.User)
	if !ok {
		http.Error(w, `{"error":"invalid user context"}`, http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"invalid input!"}`, http.StatusBadRequest)
		return
	}

	for key, values := range r.Form {
		h.logger.Info(">>> [DEBUG] Поле: %s, Значение: %v\n", key, values)
	}

	// Если передано имя — обновляем кастомное имя пользователя
	name := r.FormValue("name")
	if name != "" {
		h.logger.Info(">>> [DEBUG] Обновляем имя пользователя: %s\n", name)
		user.CustomName = name
		_, err := h.userService.UpdateUser(ctx, *user)
		if err != nil {
			h.logger.Error(">>> [ERROR] Ошибка обновления имени: %v\n", err)
			http.Error(w, `{"error":"failed to update custom name"}`, http.StatusInternalServerError)
			return
		}
	}

	// Создаем пост, используя правильные имена полей из формы
	post := entity.Post{
		Title:    r.FormValue("subject"),
		Content:  r.FormValue("comment"),
		ImageURL: "",
		UserID:   user.ID,
	}

	// Обработка файла (если он был загружен)
	file, handler, err := r.FormFile("file")
	if err == nil && handler != nil {
		defer file.Close()

		// Считываем размер файла
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, `{"error":"failed to read file"}`, http.StatusInternalServerError)
			return
		}
		fileReader := strings.NewReader(string(fileBytes))

		objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)

		url, err := h.fileStorage.UploadFile(ctx, "posts", objectName, fileReader, int64(len(fileBytes)), handler.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, `{"error":"failed to upload file"}`, http.StatusInternalServerError)
			return
		}

		post.ImageURL = url
	}

	h.logger.Info(">>> [DEBUG] Создаем пост: Title=%s, Content=%s, ImageURL=%s, UserID=%d\n",
		post.Title, post.Content, post.ImageURL, post.UserID)

	createdPost, err := h.postService.CreatePost(ctx, post)
	if err != nil {
		h.logger.Error(">>> [ERROR] Ошибка создания поста: %v\n", err)
		http.Error(w, `{"error":"failed to create post"}`, http.StatusInternalServerError)
		return
	}

	h.logger.Info(">>> [DEBUG] Пост успешно создан: ID=%d\n", createdPost.ID)

	if err := json.NewEncoder(w).Encode(createdPost); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>>GetPostByID handler called")
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
		h.logger.Warn("Dont have comments or error:", err)
		comments = []entity.Comment{}
	}

	cmVM := make([]CommentVM, 0, len(comments))
	for _, c := range comments {
		if c == (entity.Comment{}) {
			continue
		}
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
			ParentID:  c.ParentID,
			Image:     c.FileURL,
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
		Deleted:    post.DeletedAt,
	}

	// 7) Рендерим шаблон
	for _, c := range comments {
		h.logger.Info("CommentID: %d, ParentID: %d\n", c.ID, c.ParentID)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := RenderTemplate(w, "post.html", vm); err != nil {
		http.Error(w, "Error template rendering: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdatePostPostByID(w http.ResponseWriter, r *http.Request) {
}

func (h *PostHandler) DeletePostPostPostByID(w http.ResponseWriter, r *http.Request) {
}

func (h *PostHandler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, "create-post.html", nil)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info(">>>GetPostsHandler called")
	ctx := r.Context()
	posts, err := h.postService.ListPosts(ctx)
	if err != nil {
		h.logger.Error("Ошибка при получении постов", http.StatusInternalServerError)
		http.Error(w, `{"error":"error receiving posts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = RenderTemplate(w, "catalog.html", posts)
	if err != nil {
		http.Error(w, "Error template rendering", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) GetArchivedPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.postService == nil {
		http.Error(w, `{"error":"internal server error - service not initialized"}`, http.StatusInternalServerError)
		return
	}

	h.logger.Debug("GetArchivedPosts: начинаем выполнение")

	posts, err := h.postService.ListArchivedPosts(ctx)
	if err != nil {
		h.logger.Error("Ошибка при получении постов", err)
		http.Error(w, `{"error":"error receiving posts"}`, http.StatusInternalServerError)
		return
	}
	if posts == nil {
		posts = []entity.PostRequest{}
	}
	h.logger.Debug("Fetched archived posts:", posts)

	w.Header().Set("Content-Type", "text/html")

	tmplPath := filepath.Join("web", "archive.html")
	if _, err := template.ParseFiles(tmplPath); err != nil {
		h.logger.Error("Ошибка при парсинге шаблона", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = RenderTemplate(w, "archive.html", posts)
	if err != nil {
		h.logger.Error("Ошибка при рендеринге шаблона", err)
		http.Error(w, "Error template rendering", http.StatusInternalServerError)
		return
	}
}
