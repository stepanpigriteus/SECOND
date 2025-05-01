package http

import (
	"a1337b04rd/internal/adapter/http/handler"
	"a1337b04rd/internal/adapter/http/router"
	"a1337b04rd/internal/adapter/repository/postgress"
	"a1337b04rd/internal/api"
	"a1337b04rd/internal/domain/port"
	"a1337b04rd/internal/domain/service"
	"a1337b04rd/pkg/errors"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	servicework "a1337b04rd/internal/domain/service/service_work"
)

type server struct {
	port     string
	db       *sql.DB
	logger   port.Logger
	services *service.AllServices
}

func NewServer(port string, db *sql.DB, logger port.Logger) *server {
	minioStorage, err := api.NewMinioStorage("minio:9000", "minioadmin_new", "minioadmin_new_password", "posts", false)
	if err != nil {
		log.Fatal(err)
	}

	fileStorageService := servicework.NewFileStorageService(minioStorage)

	// Инициализация репозиториев и сервисов
	postRepo := postgress.NewPostgresPostRepository(db)
	postService := servicework.NewPostService(postRepo)
	userRepo := postgress.NewPostgresUserRepository(db)
	userService := servicework.NewUserService(userRepo)
	commentRepo := postgress.NewPostgresCommentRepository(db)
	commentService := servicework.NewCommentService(commentRepo)
	sessionRepo := postgress.NewPostgresSessionRepository(db)
	sessionService := servicework.NewSessionService(sessionRepo)

	services := &service.AllServices{
		Post:    postService,
		User:    userService,
		Comment: commentService,
		Session: sessionService,
		Storage: fileStorageService,
	}

	return &server{
		port:     port,
		db:       db,
		logger:   logger,
		services: services,
	}
}

type handleDef struct{}

func (s *server) RunServer() {
	if s.port == "" {
		s.logger.Error("Port is not set")
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Инициализация хендлеров
	handlers := handler.NewAllHandlers(s.services)

	// Регистрируем роуты
	router.RegisterRoutes(mux, handlers, s.services.Session, s.services.User)

	s.logger.Info("Starting server", "port", s.port)

	err := http.ListenAndServe("0.0.0.0:"+s.port, mux)
	if err != nil {
		s.logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func (h *handleDef) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusInternalServerError
	if r.Method == "OPTIONS" {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)

	response := errors.Error{
		Message: "Undefined Error, please check your method or endpoint correctness",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}
