package store

import (
	"context"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStoreTasks struct {
	tasks *mongo.Collection
}

func (t MongoDBStoreTasks) Insert(ctx context.Context, task *data.Task) error {
	task.ID = primitive.NewObjectID()

	_, err := t.tasks.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	return nil
}
func (t MongoDBStoreTasks) GetTaskById(context.Context, string) (*data.Task, error) {
	return nil, nil
}
func (t MongoDBStoreTasks) GetOwnerOfTask(context.Context, string) (string, error) {
	return "", nil
}
func (t MongoDBStoreTasks) FetchAllTasks(context.Context, string, []string) ([]*data.Task, error) {
	return nil, nil
}
func (t MongoDBStoreTasks) UpdateTask(context.Context, string, *data.Task) (*data.Task, error) {
	return nil, nil
}
func (t MongoDBStoreTasks) DeleteTask(context.Context, string) error {
	return nil
}
func (t MongoDBStoreTasks) CompleteTask(context.Context, string) error {
	return nil
}
