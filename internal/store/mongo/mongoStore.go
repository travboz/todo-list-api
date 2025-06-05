package mongo

import (
	"context"

	"github.com/travboz/backend-projects/todo-list-api/internal/env"
	"github.com/travboz/backend-projects/todo-list-api/internal/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBStorage(db *mongo.Client) (*store.Storage, error) {
	dbName := env.GetString("MONGO_DB_NAME", "todo-list-api")

	_, err := db.Database(dbName).Collection("tokens").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(360), // 6 minute expiry
	})

	if err != nil {
		return nil, err
	}

	_, err = db.Database(dbName).Collection("tasks").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"title", "text"}},
	})

	if err != nil {
		return nil, err
	}

	return &store.Storage{
		UsersModel:  MongoDBStoreUsers{Users: db.Database(dbName).Collection("users")},
		TasksModel:  MongoDBStoreTasks{Tasks: db.Database(dbName).Collection("tasks")},
		TokensModel: MongoDBStoreTokens{Tokens: db.Database(dbName).Collection("tokens")},
	}, nil
}
