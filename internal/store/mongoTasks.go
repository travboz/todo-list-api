package store

import (
	"context"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStoreTasks struct {
	tasks *mongo.Collection
}

func (t MongoStoreTasks) Insert(context.Context, *data.Task) error {
	return nil
}
func (t MongoStoreTasks) GetTaskById(context.Context, string) (*data.Task, error) {
	return nil, nil
}
func (t MongoStoreTasks) GetOwnerOfTask(context.Context, string) (string, error) {
	return "", nil
}
func (t MongoStoreTasks) FetchAllTasks(context.Context, string, []string) ([]*data.Task, error) {
	return nil, nil
}
func (t MongoStoreTasks) UpdateTask(context.Context, string, *data.Task) (*data.Task, error) {
	return nil, nil
}
func (t MongoStoreTasks) DeleteTask(context.Context, string) error {
	return nil
}
func (t MongoStoreTasks) CompleteTask(context.Context, string) error {
	return nil
}
