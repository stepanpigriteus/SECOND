package http

import (
	"1337b04rd/internal/domain/port"
	"database/sql"
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

func (s *server) RunServer() {
	
	os.Exit(1)
}
