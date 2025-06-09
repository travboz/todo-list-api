package store

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type TasksStore struct {
	db    *mongo.Collection
	cache *redis.Client
}

// NewTasksStore creates a new TasksStore that implements TasksModel interface
func NewTasksStore(collection *mongo.Collection, cache *redis.Client) TasksModel {
	return &TasksStore{
		db:    collection,
		cache: cache,
	}
}

func (ts *TasksStore) Insert(context.Context, *data.Task) error {
	return nil
}

func (ts *TasksStore) GetTaskById(context.Context, string) (*data.Task, error) {
	return nil, nil
}

func (ts *TasksStore) FetchAllTasks(ctx context.Context, p data.Filters, search string) ([]*data.Task, data.Metadata, error) {
	return nil, data.Metadata{}, nil
}

func (ts *TasksStore) UpdateTask(ctx context.Context, id string, t *data.Task) (*data.Task, error) {
	return nil, nil
}

func (ts *TasksStore) DeleteTask(context.Context, string) error {
	return nil
}

func (ts *TasksStore) CompleteTask(ctx context.Context, id string) (*data.Task, error) {
	return nil, nil
}

func (ts *TasksStore) GetTaskOwnerId(ctx context.Context, id string) (string, error) {
	return "", nil
}
