package main

import (
	"fmt"
	"log"

	"migrations/internal/config"
	"migrations/internal/server"
	"migrations/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	db, err := store.NewDB(cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to initialize a new DB: %w", err)
	}
	h := server.NewHandlers(db)
	return server.RunServer(h)
}
