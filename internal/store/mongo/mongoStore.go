package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client *mongo.Client
	dbName string
	db     *mongo.Database
}

func NewMongoStore(client *mongo.Client, dbName string) (*store.Storage, error) {
	ms := &MongoStorage{
		client: client,
		dbName: dbName,
		db:     client.Database(dbName),
	}

	// Setup all indexes
	if err := ms.setupIndexes(); err != nil {
		return nil, fmt.Errorf("failed to setup indexes: %w", err)
	}

	// Create and return storage with all models
	return &store.Storage{
		UsersModel:  ms.NewMongoUsersModel(),
		TasksModel:  ms.NewMongoTasksModel(),
		TokensModel: ms.NewMongoTokensModel(),
	}, nil
}

func (ms *MongoStorage) setupIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup tokens index with TTL
	if err := ms.setupTokensIndexes(ctx); err != nil {
		return fmt.Errorf("tokens indexes: %w", err)
	}

	// Setup tasks search index
	if err := ms.setupTasksIndexes(ctx); err != nil {
		return fmt.Errorf("tasks indexes: %w", err)
	}

	return nil
}

func (ms *MongoStorage) setupTokensIndexes(ctx context.Context) error {
	collection := ms.db.Collection("tokens")

	// TTL index for token expiry
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(360), // 6 minutes
	})

	return err
}

func (ms *MongoStorage) setupTasksIndexes(ctx context.Context) error {
	collection := ms.db.Collection("tasks")

	// Text search index
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "title", Value: "text"}},
	})

	return err
}

// func (ms *MongoStorage) setupUsersIndexes(ctx context.Context) error {
// 	collection := ms.db.Collection("users")

// 	// Email uniqueness index
// 	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
// 		Keys:    bson.D{{Key: "email", Value: 1}},
// 		Options: options.Index().SetUnique(true),
// 	})

// 	return err
// }
