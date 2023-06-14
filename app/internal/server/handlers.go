package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"migrations/internal/model"
	"migrations/internal/store"
)

type Handlers struct {
	store *store.DB
}

func NewHandlers(s *store.DB) *Handlers {
	return &Handlers{
		store: s,
	}
}

func (h *Handlers) PutEmployee(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read the PutEmployee request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var emp model.Employee
	if err := json.Unmarshal(b, &emp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.store.PutEmployee(ctx, &emp); err != nil {
		log.Printf("failed to store employee in the DB: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
