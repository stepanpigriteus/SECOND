package main

import (
	"1337b04rd/internal/adapter/http"
	repository "1337b04rd/internal/adapter/repository/postgress"
	"1337b04rd/internal/config"
	"1337b04rd/pkg/flags"
	"1337b04rd/pkg/logger"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	port := flags.Flags()
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	logger := logger.NewSlogAdapter()
	db, err := repository.ConnectDB(connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	server := http.NewServer(strconv.Itoa(port), db, logger)
	server.RunServer()
}
