package store

import (
	"context"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"github.com/travboz/backend-projects/todo-list-api/internal/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	UsersModel
	TasksModel
	TokensModel
}

func NewMongoDBStorage(db *mongo.Client) *Storage {
	dbName := env.GetString("MONGO_DB_NAME", "todo-list-api")

	db.Database(dbName).Collection("tokens").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(360), // 6 minute expiry
	})

	return &Storage{
		UsersModel:  MongoDBStoreUsers{db.Database(dbName).Collection("users")},
		TasksModel:  MongoDBStoreTasks{db.Database(dbName).Collection("tasks")},
		TokensModel: MongoDBStoreTokens{db.Database(dbName).Collection("tokens")},
	}
}

type TasksModel interface {
	Insert(context.Context, *data.Task) error
	GetTaskById(context.Context, string) (*data.Task, error)
	FetchAllTasks(ctx context.Context) ([]*data.Task, error)
	UpdateTask(context.Context, string, *data.Task) (*data.Task, error)
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
