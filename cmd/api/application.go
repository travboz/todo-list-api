package main

import (
	"log/slog"

	"github.com/travboz/backend-projects/todo-list-api/internal/store"
)

type application struct {
	Logger  *slog.Logger
	Storage *store.Storage
}
