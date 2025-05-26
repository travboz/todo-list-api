package store

import (
	"context"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
)

type Store struct {
	UsersModel
	TasksModel
}

func NewStore(t TasksModel, u UsersModel) *Store {
	return &Store{
		UsersModel: u,
		TasksModel: t,
	}
}

type TasksModel interface {
	Insert(context.Context, *data.Task) error
	GetTaskById(context.Context, string) (*data.Task, error)
	// GetOwnerOfTask(context.Context, string) (string, error)
	// FetchAllTasks(context.Context, string, []string) ([]*data.Task, error)
	// UpdateTask(context.Context, string, *data.Task) (*data.Task, error)
	// DeleteTask(context.Context, string) error
	// CompleteTask(context.Context, string) error
	// Shutdown(context.Context) error
}

type UsersModel interface {
	Insert(name, email, password string) error
	// Exists(id string) (bool, error)
	// Get(id string) (*data.User, error)
}
