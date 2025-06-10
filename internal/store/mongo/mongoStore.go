package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Client *mongo.Client
	DbName string
	DB     *mongo.Database
}

func NewMongoStore(client *mongo.Client, dbName string) (*MongoStorage, error) {
	ms := &MongoStorage{
		Client: client,
		DbName: dbName,
		DB:     client.Database(dbName),
	}

	// Setup all indexes
	if err := setupIndexes(ms.DB); err != nil {
		return nil, fmt.Errorf("failed to setup indexes: %w", err)
	}

	return ms, nil
}

func setupIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup tokens index with TTL
	if err := setupTokensIndexes(ctx, db); err != nil {
		return fmt.Errorf("tokens indexes: %w", err)
	}

	// Setup tasks search index
	if err := setupTasksIndexes(ctx, db); err != nil {
		return fmt.Errorf("tasks indexes: %w", err)
	}

	return nil
}

func setupTokensIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("tokens")

	// TTL index for token expiry
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(360), // 6 minutes
	})

	return err
}

func setupTasksIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("tasks")

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
