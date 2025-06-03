package store

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	tokenGen "github.com/travboz/backend-projects/todo-list-api/internal/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStoreTokens struct {
	Tokens *mongo.Collection
}

func (m MongoDBStoreTokens) InsertToken(ctx context.Context, user_id string) (string, error) {
	var token data.Token

	token.ID = primitive.NewObjectID()
	token.CreatedAt = time.Now()
	token.Token = tokenGen.GenerateToken()
	token.UserID, _ = primitive.ObjectIDFromHex(user_id)

	_, err := m.Tokens.InsertOne(ctx, token)

	return token.Token, err
}

func (m MongoDBStoreTokens) ValidateToken(ctx context.Context, token string) (bool, error) {
	filter := bson.M{"token": token}
	result := m.Tokens.FindOne(ctx, filter)

	var tokenData data.Token
	if err := result.Decode(&tokenData); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return false, ErrRecordNotFound
		default:
			return false, err
		}
	}

	if subtle.ConstantTimeCompare([]byte(token), []byte(tokenData.Token)) == 1 {
		return true, nil
	}

	return false, nil
}

func (m MongoDBStoreTokens) GetUserIdUsingToken(ctx context.Context, token string) (string, error) {
	filter := bson.M{"token": token}
	result := m.Tokens.FindOne(ctx, filter)

	var tokenData data.Token
	if err := result.Decode(&tokenData); err != nil {
		return "", ErrRecordNotFound
	}

	return tokenData.UserID.Hex(), nil
}
