package mongo

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	appErrors "github.com/travboz/backend-projects/todo-list-api/internal/errors"
	"github.com/travboz/backend-projects/todo-list-api/internal/store"
	"github.com/travboz/backend-projects/todo-list-api/internal/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStoreTokens struct {
	Tokens *mongo.Collection
}

func (ms *MongoStorage) NewMongoTokensModel() store.TokensModel {
	return MongoDBStoreTokens{ms.DB.Collection("tokens")}
}

func (m MongoDBStoreTokens) InsertToken(ctx context.Context, user_id string) (string, error) {
	var newToken data.Token

	rand_token, err := token.GenerateToken(32)
	if err != nil {
		return "", err
	}

	newToken.ID = primitive.NewObjectID()
	newToken.CreatedAt = time.Now()
	newToken.Token = rand_token

	newToken.UserID, err = primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return "", err
	}

	_, err = m.Tokens.InsertOne(ctx, newToken)

	return newToken.Token, err
}

// ValidateToken looks up the token in Mongo to see if it still exists. As each token entry has an expiry time.
// It then compares that the token matches exactly (just in case).
func (m MongoDBStoreTokens) ValidateToken(ctx context.Context, token string) (bool, error) {
	filter := bson.M{"token": token}
	result := m.Tokens.FindOne(ctx, filter)

	var tokenData data.Token
	if err := result.Decode(&tokenData); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return false, appErrors.ErrRecordNotFound
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
		return "", appErrors.ErrRecordNotFound
	}

	return tokenData.UserID.Hex(), nil
}
