package http

import (
	"1337b04rd/internal/adapter/http/handler"
	"1337b04rd/internal/adapter/http/router"
	"1337b04rd/internal/adapter/repository/postgress"
	"1337b04rd/internal/domain/port"
	"1337b04rd/internal/domain/service"
	servicework "1337b04rd/internal/domain/service/service_work"
	"1337b04rd/pkg/errors"
	"1337b04rd/pkg/logger"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
)

type server struct {
	port     string
	db       *sql.DB
	logger   port.Logger
	services *service.AllServices
}

func NewServer(port string, db *sql.DB, logger port.Logger) *server {
	// Инициализация репозиториев и сервисов
	postRepo := postgress.NewPostgresPostRepository(db)
	postService := servicework.NewPostService(postRepo)

	services := &service.AllServices{
		Post: postService,
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
		logger.NewSlogAdapter().Error("Port is not set")
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Инициализация хендлеров
	handlers := handler.NewAllHandlers(s.services)

	// Регистрируем роуты
	router.RegisterRoutes(mux, handlers)

	logger.NewSlogAdapter().Info("Starting server", "port", s.port)

	err := http.ListenAndServe("0.0.0.0:"+s.port, mux)
	if err != nil {
		logger.NewSlogAdapter().Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func (h *handleDef) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	response := errors.Error{Message: "Undefined Error, please check your method or endpoint correctness"}
	json.NewEncoder(w).Encode(response)
}
