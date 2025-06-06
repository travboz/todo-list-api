package store

import (
	"context"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
)

type Storage struct {
	UsersModel
	TasksModel
	TokensModel
}

type TasksModel interface {
	Insert(context.Context, *data.Task) error
	GetTaskById(context.Context, string) (*data.Task, error)
	FetchAllTasks(ctx context.Context, p data.Filters, search string) ([]*data.Task, data.Metadata, error)
	UpdateTask(ctx context.Context, id string, t *data.Task) (*data.Task, error)
	DeleteTask(context.Context, string) error
	CompleteTask(ctx context.Context, id string) (*data.Task, error)
	GetTaskOwnerId(ctx context.Context, id string) (string, error)

	// Shutdown(context.Context) error
}

type UsersModel interface {
	Insert(u *data.User) error
	Authenticate(email, password string) (string, error)
	Get(id string) (*data.User, error)
}

type TokensModel interface {
	InsertToken(ctx context.Context, user_id string) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
	GetUserIdUsingToken(ctx context.Context, token string) (string, error)
}
