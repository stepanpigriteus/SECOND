package handler

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/port"
	"1337b04rd/internal/domain/port/storage"
	"1337b04rd/internal/domain/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostHandler struct {
	postService    service.PostService
	userService    service.UserService
	commentService service.CommentService
	logger         port.Logger
	fileStorage    storage.FileStorage
}

func NewPostHandler(ps service.PostService, us service.UserService, cs service.CommentService, fs storage.FileStorage) *PostHandler {
	return &PostHandler{
		postService:    ps,
		userService:    us,
		commentService: cs,
		fileStorage:    fs,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> CreatePost handler called")
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

	fmt.Println(">>> [DEBUG] Полученные данные формы:")
	for key, values := range r.Form {
		fmt.Printf(">>> [DEBUG] Поле: %s, Значение: %v\n", key, values)
	}

	// Если передано имя — обновляем кастомное имя пользователя
	name := r.FormValue("name")
	if name != "" {
		fmt.Printf(">>> [DEBUG] Обновляем имя пользователя: %s\n", name)
		user.CustomName = name
		_, err := h.userService.UpdateUser(ctx, *user)
		if err != nil {
			fmt.Printf(">>> [ERROR] Ошибка обновления имени: %v\n", err)
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

		url, err := h.fileStorage.UploadFile(ctx, objectName, fileReader, int64(len(fileBytes)), handler.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, `{"error":"failed to upload file"}`, http.StatusInternalServerError)
			return
		}

		post.ImageURL = url
	}

	fmt.Printf(">>> [DEBUG] Создаем пост: Title=%s, Content=%s, ImageURL=%s, UserID=%d\n",
		post.Title, post.Content, post.ImageURL, post.UserID)

	createdPost, err := h.postService.CreatePost(ctx, post)
	if err != nil {
		fmt.Printf(">>> [ERROR] Ошибка создания поста: %v\n", err)
		http.Error(w, `{"error":"failed to create post"}`, http.StatusInternalServerError)
		return
	}

	fmt.Printf(">>> [DEBUG] Пост успешно создан: ID=%d\n", createdPost.ID)

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
		fmt.Println("Dont have comments or error:", err)
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
	for _, c := range comments {
		fmt.Printf("CommentID: %d, ParentID: %d\n", c.ID, c.ParentID)
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

func (h *PostHandler) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, "create-post.html", nil)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
