package store

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	appErrors "github.com/travboz/backend-projects/todo-list-api/internal/errors"
	"github.com/travboz/backend-projects/todo-list-api/internal/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokensStore struct {
	db    *mongo.Collection
	cache *redis.Client
}

// NewTasksStore creates a new TasksStore that implements TasksModel interface
func NewTokensStore(collection *mongo.Collection, cache *redis.Client) TokensModel {
	return &TokensStore{
		db:    collection,
		cache: cache,
	}
}

func (ts *TokensStore) InsertToken(ctx context.Context, user_id string) (string, error) {
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

	_, err = ts.db.InsertOne(ctx, newToken)

	return newToken.Token, err
}

func (ts *TokensStore) getTokenFromCache(ctx context.Context, key string) (*data.Token, bool) {
	val, err := ts.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, false // cache miss or error
	}

	var token data.Token
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return nil, false // failed to parse cached value
	}

	return &token, true
}

func (ts *TokensStore) isTokenEqual(known, test []byte) bool {
	return subtle.ConstantTimeCompare(known, test) == 1
}

func (ts *TokensStore) getTokenFromDB(ctx context.Context, token string) (*data.Token, error) {
	filter := bson.M{"token": token}
	result := ts.db.FindOne(ctx, filter)

	var tokenData data.Token
	if err := result.Decode(&tokenData); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, appErrors.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tokenData, nil
}

// ValidateToken looks up the token in Mongo to see if it still exists. As each token entry has an expiry time.
// It then compares that the token matches exactly (just in case).
func (ts *TokensStore) GetAndValidateToken(ctx context.Context, token string) (bool, error) {
	cacheKey := fmt.Sprintf("auth:user:%s", token)

	// 1. Try cache first
	if tok, ok := ts.getTokenFromCache(ctx, cacheKey); ok {
		if ts.isTokenEqual([]byte(token), []byte(tok.Token)) {
			return true, nil
		}
	}

	// 2. Cache miss or error - try database
	tok, err := ts.getTokenFromDB(ctx, token)
	if err != nil {
		return false, err
	}

	// 3. Found in database - cache it for next time
	if ts.isTokenEqual([]byte(token), []byte(tok.Token)) {
		ts.cacheToken(ctx, cacheKey, tok)

		return true, nil
	}

	return false, nil
}

func (ts *TokensStore) cacheToken(ctx context.Context, key string, token *data.Token) {
	if data, err := json.Marshal(token); err == nil {
		ts.cache.Set(ctx, key, data, CacheExpiryTime)
	}
}

func (ts *TokensStore) GetUserIdUsingToken(ctx context.Context, token string) (string, error) {
	filter := bson.M{"token": token}
	result := ts.db.FindOne(ctx, filter)

	var tokenData data.Token
	if err := result.Decode(&tokenData); err != nil {
		return "", appErrors.ErrRecordNotFound
	}

	return tokenData.UserID.Hex(), nil
}
