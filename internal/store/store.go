package store

import (
	"context"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"github.com/travboz/backend-projects/todo-list-api/internal/env"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	UsersModel
	TasksModel
}

func NewMongoDBStorage(db *mongo.Client) *Storage {
	dbName := env.GetString("MONGO_DB_NAME", "todo-list-api")

	return &Storage{
		UsersModel: MongoDBStoreUsers{db.Database(dbName).Collection("users")},
		TasksModel: MongoDBStoreTasks{db.Database(dbName).Collection("tasks")},
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
