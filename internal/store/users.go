package store

import (
	"context"
	"errors"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type MongoDBStoreUsers struct {
	Users *mongo.Collection
}

func (u MongoDBStoreUsers) Insert(user *data.User) error {
	user.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	_, err = u.Users.InsertOne(ctx, user)

	return err
}

func (u MongoDBStoreUsers) Authenticate(email, password string) (string, error) {
	filter := bson.M{"email": email}

	result := u.Users.FindOne(context.Background(), filter)

	var user data.User
	if err := result.Decode(&user); err != nil {
		var user data.User
		if err = result.Decode(&user); err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				return "", ErrRecordNotFound
			default:
				return "", err
			}
		}
	}

	// check whether passwords match
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}

	}

	// if password is correct, return user ID
	return user.ID.Hex(), nil

}

func (u MongoDBStoreUsers) Get(id string) (*data.User, error) {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": user_id}

	result := u.Users.FindOne(ctx, filter)

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
