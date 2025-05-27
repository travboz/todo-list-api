package store

import (
	"context"
	"errors"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStoreUsers struct {
	users *mongo.Collection
}

func (u MongoDBStoreUsers) Insert(user *data.User) error {
	user.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := u.users.InsertOne(ctx, user)

	return err
}

func (u MongoDBStoreUsers) Exists(id string) (bool, error) {
	return false, nil
}

func (u MongoDBStoreUsers) Get(id string) (*data.User, error) {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": user_id}

	result := u.users.FindOne(ctx, filter)

	var user data.User
	if err = result.Decode(&user); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
