package main

import (
	"1337b04rd/pkg/flags"
	"1337b04rd/pkg/logger"
)

func main() {
	port := flags.Flags()

	logger := logger.NewSlogAdapter()
	logger.Info("User created", "user_id", port, "email", "rick@example.com")
}
