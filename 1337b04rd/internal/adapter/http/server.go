package http

import (
	"1337b04rd/internal/domain/port"
	"1337b04rd/pkg/errors"
	"1337b04rd/pkg/logger"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
)

type server struct {
	port   string
	db     *sql.DB
	logger *port.Logger
}

func NewServer(port string, db *sql.DB, logger *port.Logger) *server {
	return &server{
		port:   port,
		db:     db,
		logger: logger,
	}
}

type handleDef struct{}

func (s *server) RunServer() {
	if s.port == "" {
		logger.NewSlogAdapter().Error("Port is not set")
		os.Exit(1)
	}
	mux := http.NewServeMux()
	mux.Handle("/", &handleDef{})
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
